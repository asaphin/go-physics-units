package units

import (
	"cmp"
	"errors"
	"maps"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/asaphin/go-physics-units/internal/rates"
)

// Parser parses measurements from strings, like ParseString, with options that
// ParseString doesn't have: unit aliases, case-insensitive symbols, the form
// without a space ("10km") and user-defined measures. Build one with NewParser.
//
// Like ParseString, Parse returns values of the measure's concrete type: a
// Distance for "10 km", a Quantity[M] for a unit of a measure added with
// WithMeasure.
//
// A Parser is immutable and safe for concurrent use. It is a small value that
// refers to data no one can change, so copy it freely. The zero Parser parses
// exactly like ParseString.
type Parser struct {
	p *parser
}

type parser struct {
	units   map[string]parsedUnit // canonical symbols and aliases
	folded  map[string]parsedUnit // lower-case symbols and aliases; nil unless CaseInsensitive
	noSpace bool
}

// ParserOption configures a Parser; see NewParser.
type ParserOption func(*parserConfig)

type parserConfig struct {
	aliases         []map[string]string
	measures        []parsedMeasure
	caseInsensitive bool
	noSpace         bool
	err             error
}

// WithAliases adds alternative names for units. Each key parses as the unit
// symbol it maps to:
//
//	units.WithAliases(map[string]string{"metre": distance.Meter, "kph": velocity.KilometerPerHour})
//
// The targets must be predefined symbols, or units of a measure added with
// WithMeasure. An alias must be non-empty, contain no whitespace, and not be a
// symbol or another alias. The map is copied.
func WithAliases(aliases map[string]string) ParserOption {
	copied := maps.Clone(aliases)

	return func(c *parserConfig) {
		c.aliases = append(c.aliases, copied)
	}
}

// CaseInsensitive makes symbols and aliases match regardless of case where
// that is unambiguous: "10 KM" parses as 10 km. An exact match always wins, so
// "1 kn" is still a knot, but "1 KN" fails with ErrAmbiguousUnit, because kN
// (kilonewton) and kn (knot) differ only in case.
func CaseInsensitive() ParserOption {
	return func(c *parserConfig) {
		c.caseInsensitive = true
	}
}

// AllowNoSpace also accepts a value directly followed by its unit, as in
// "10km" or "1.5e3eV". In that form the unit can't start with a digit, a sign
// or a point. Strings with whitespace between the two still parse.
func AllowNoSpace() ParserOption {
	return func(c *parserConfig) {
		c.noSpace = true
	}
}

// parsedMeasure is a user-defined measure added with WithMeasure.
type parsedMeasure struct {
	set *unitSet
	box func(float64, uint16) Measurement
}

// WithMeasure makes the parser recognise the units of the user-defined measure
// M. Parsed values are of type Quantity[M]. The measure's symbols must not
// clash with predefined symbols or with those of another measure.
func WithMeasure[M Measure]() ParserOption {
	return func(c *parserConfig) {
		s, err := setOf[M]()
		if err != nil {
			c.err = cmp.Or(c.err, invalidOption("WithMeasure: "+err.Error()))

			return
		}

		c.measures = append(c.measures, parsedMeasure{set: s, box: func(v float64, u uint16) Measurement {
			return Quantity[M]{value: v, unit: Unit[M]{i: u}}
		}})
	}
}

// NewParser builds a Parser from the options. It fails with ErrInvalidOption if
// an alias is invalid or doesn't denote a unit, or if two measures share a
// symbol.
func NewParser(options ...ParserOption) (Parser, error) {
	var c parserConfig

	for _, option := range options {
		option(&c)
	}

	if c.err != nil {
		return Parser{}, c.err
	}

	p := &parser{units: maps.Clone(predefinedUnits), noSpace: c.noSpace}

	for _, m := range c.measures {
		if err := p.addMeasure(m); err != nil {
			return Parser{}, err
		}
	}

	for _, aliases := range c.aliases {
		if err := p.addAliases(aliases); err != nil {
			return Parser{}, err
		}
	}

	if c.caseInsensitive {
		p.fold()
	}

	return Parser{p: p}, nil
}

func (p *parser) addMeasure(m parsedMeasure) error {
	for i, symbol := range m.set.symbols {
		if prev, ok := p.units[symbol]; ok {
			return invalidOption("unit " + symbol + " of measure " + string(m.set.name) +
				" is already a unit of measure " + string(prev.set.name))
		}

		p.units[symbol] = parsedUnit{
			set:  m.set,
			unit: uint16(i),
			box:  m.box,
		}
	}

	return nil
}

func (p *parser) addAliases(aliases map[string]string) error {
	for alias, target := range aliases {
		canonical := rates.CanonicalUnit(alias)
		if canonical == "" || strings.ContainsFunc(canonical, unicode.IsSpace) {
			return invalidOption("alias " + strconv.Quote(alias) + " must be non-empty and contain no whitespace")
		}

		if _, ok := p.units[canonical]; ok {
			return invalidOption("alias " + alias + " is already a unit symbol or alias")
		}

		// The target must be a symbol, not another alias: then the result
		// doesn't depend on the order in which aliases are added.
		resolved, ok := p.units[rates.CanonicalUnit(target)]
		if !ok || resolved.set.symbols[resolved.unit] != rates.CanonicalUnit(target) {
			return invalidOption("alias " + alias + ": " + target + " isn't a unit symbol")
		}

		p.units[canonical] = resolved
	}

	return nil
}

// fold builds the case-insensitive index. A lower-case key that several
// symbols or aliases with different meanings share is ambiguous.
func (p *parser) fold() {
	groups := make(map[string][]string, len(p.units))
	for key := range p.units {
		lower := strings.ToLower(key)
		groups[lower] = append(groups[lower], key)
	}

	p.folded = make(map[string]parsedUnit, len(groups))

	for lower, keys := range groups {
		first := p.units[keys[0]]

		ambiguous := slices.ContainsFunc(keys[1:], func(k string) bool {
			u := p.units[k]

			return u.unit != first.unit || u.set != first.set
		})
		if !ambiguous {
			p.folded[lower] = first

			continue
		}

		slices.Sort(keys)
		p.folded[lower] = parsedUnit{candidates: strings.Join(keys, ", ")}
	}
}

// Parse parses a measurement from a string of the form "<value> <unit>", where
// the unit may be any symbol or alias the parser accepts. The value must be
// finite. Values of predefined measures have their measure's type, such as
// Distance, and values of a measure M added with WithMeasure the type
// Quantity[M].
func (p Parser) Parse(s string) (Measurement, error) {
	if p.p == nil {
		return ParseString(s)
	}

	value, symbol, err := p.p.split(s)
	if err != nil {
		return nil, err
	}

	resolved, err := p.p.resolve(symbol)
	if err != nil {
		return nil, err
	}

	return resolved.box(value, resolved.unit), nil
}

// Unit returns the canonical unit symbol that symbol denotes: itself, or the
// symbol an alias or a case variant stands for. Use it with the unit
// packages' ParseUnit:
//
//	symbol, err := p.Unit("metres") // "m"
//	u, ok := distance.ParseUnit(symbol)
func (p Parser) Unit(symbol string) (string, error) {
	if p.p == nil {
		pu, ok := predefinedUnits[rates.CanonicalUnit(symbol)]
		if !ok {
			return "", unknownUnitError("", symbol)
		}

		return pu.set.symbols[pu.unit], nil
	}

	resolved, err := p.p.resolve(symbol)
	if err != nil {
		return "", err
	}

	return resolved.set.symbols[resolved.unit], nil
}

func (p *parser) resolve(symbol string) (parsedUnit, error) {
	canonical := rates.CanonicalUnit(symbol)

	if resolved, ok := p.units[canonical]; ok {
		return resolved, nil
	}

	if p.folded != nil {
		if resolved, ok := p.lookupFolded(canonical); ok {
			if resolved.candidates != "" {
				return parsedUnit{}, newError(ErrAmbiguousUnit,
					"ambiguous unit "+symbol+": it could be any of "+resolved.candidates, nil)
			}

			return resolved, nil
		}
	}

	return parsedUnit{}, unknownUnitError("", symbol)
}

// lookupFolded looks canonical up in the case-insensitive index. Short ASCII
// symbols, the common case, are lowered into a stack buffer: a map lookup with
// string(bytes) doesn't allocate, while strings.ToLower would.
func (p *parser) lookupFolded(canonical string) (parsedUnit, bool) {
	var buf [32]byte

	if len(canonical) > len(buf) || strings.ContainsFunc(canonical, func(r rune) bool { return r >= utf8.RuneSelf }) {
		resolved, ok := p.folded[strings.ToLower(canonical)]

		return resolved, ok
	}

	lower := buf[:len(canonical)]
	for i := range len(canonical) {
		c := canonical[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}

		lower[i] = c
	}

	resolved, ok := p.folded[string(lower)]

	return resolved, ok
}

// split separates the value from the unit, in the form "<value> <unit>" or,
// with AllowNoSpace, "<value><unit>".
func (p *parser) split(s string) (float64, string, error) {
	field := strings.TrimSpace(s)
	if !p.noSpace || field == "" || strings.ContainsFunc(field, unicode.IsSpace) {
		return parseValueAndUnit(s)
	}

	value, i, err := splitNumber(field)
	if err != nil {
		return 0, "", err
	}

	if i == 0 {
		return 0, "", newError(ErrInvalidFormat, "invalid string format: "+s, nil)
	}

	// The value is finite: only number characters reach ParseFloat, so it
	// can't read NaN or Inf, and an overflow is an error.
	return value, field[i:], nil
}

// splitNumber finds the value at the start of field in the form without a
// space: the longest prefix of number characters that parses as a float and
// leaves a unit. "1.5e3eV" is 1.5e3 and "eV"; "10eV" is 10 and "eV", because
// "10e" doesn't parse. It returns the value and the index where the unit
// starts, or 0 if there is no such prefix.
func splitNumber(field string) (float64, int, error) {
	end := 0
	for end < len(field) && strings.IndexByte("0123456789+-.eE", field[end]) >= 0 {
		end++
	}

	for i := end; i > 0; i-- {
		// A unit can't start with a digit, a sign or a point: that would
		// split the number itself ("10" isn't 1 of a unit "0").
		if i == len(field) || strings.IndexByte("0123456789+-.", field[i]) >= 0 {
			continue
		}

		value, err := strconv.ParseFloat(field[:i], 64)
		if errors.Is(err, strconv.ErrRange) {
			return 0, 0, newError(ErrInvalidFormat, "failed to parse value "+field[:i]+": "+err.Error(), err)
		}

		if err == nil {
			return value, i, nil
		}
	}

	return 0, 0, nil
}

// invalidOption returns an ErrInvalidOption with the given message.
func invalidOption(msg string) error {
	return newError(ErrInvalidOption, msg, nil)
}
