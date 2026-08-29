package runtime

import (
	"fmt"
	"sort"
	"strings"
)

func init() {
	for _, offset := range []string{"0", "1", "2"} {
		registerLoweringKernel("do_duplicated", offset, kernelMatchFamily)
	}
	for _, entry := range []string{"do_match", "do_pmatch", "do_makeunique", "do_vhash"} {
		registerLoweringKernel(entry, "0", kernelMatchFamily)
	}
	for _, entry := range []string{"do_xtfrm", "do_sort", "do_isunsorted", "do_sorted_fpass", "do_psort", "do_order", "do_rank"} {
		registerLoweringKernel(entry, "0", kernelOrderFamily)
	}
	registerLoweringKernel("do_sort", "1", kernelOrderFamily)
}

func vectorKeys(v Value) []string {
	keys := make([]string, Length(v))
	for i := range keys {
		keys[i] = string(v.Kind()) + ":" + takePositions(v, []int{i}).String()
	}
	return keys
}

func kernelMatchFamily(c *Context, frame *LoweringFrame) error {
	x, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	keys := vectorKeys(x)
	switch frame.Plan.CEntry {
	case "do_duplicated":
		seen := map[string]bool{}
		duplicate := make([]Logical, len(keys))
		positions := make([]int, 0, len(keys))
		first := 0
		for i, key := range keys {
			duplicate[i] = False
			if seen[key] {
				duplicate[i] = True
				if first == 0 {
					first = i + 1
				}
			} else {
				seen[key] = true
				positions = append(positions, i)
			}
		}
		switch frame.Plan.Offset {
		case "0":
			frame.Result = &LogicalVector{Data: duplicate}
		case "1":
			frame.Result = takePositions(x, positions)
		default:
			frame.Result = &IntegerVector{Data: []int64{int64(first)}}
		}
		return nil
	case "do_match", "do_pmatch":
		if len(frame.Arguments) < 2 {
			return fmt.Errorf("%s expects x and table", frame.Plan.Name)
		}
		table, err := frameValue(c, frame, 1)
		if err != nil {
			return err
		}
		tableKeys := vectorKeys(table)
		lookup := map[string]int{}
		for i, key := range tableKeys {
			if _, ok := lookup[key]; !ok {
				lookup[key] = i + 1
			}
		}
		out := &IntegerVector{Data: make([]int64, len(keys)), Missing: make([]bool, len(keys))}
		nomatch, hasNomatch := int64(0), false
		if len(frame.Arguments) > 2 {
			if value, e := frameValue(c, frame, 2); e == nil {
				if n, e := scalarInt(value); e == nil {
					nomatch, hasNomatch = int64(n), true
				}
			}
		}
		incomparable := map[string]bool{}
		if len(frame.Arguments) > 3 {
			if value, e := frameValue(c, frame, 3); e == nil {
				for _, key := range vectorKeys(value) {
					incomparable[key] = true
				}
			}
		}
		for i, key := range keys {
			pos, ok := lookup[key]
			if incomparable[key] {
				ok = false
			}
			if !ok && frame.Plan.CEntry == "do_pmatch" {
				raw := takePositions(x, []int{i}).String()
				for j, t := range tableKeys {
					if strings.HasPrefix(t, string(table.Kind())+":"+raw) {
						pos = j + 1
						ok = true
						break
					}
				}
			}
			if ok {
				out.Data[i] = int64(pos)
			} else if hasNomatch {
				out.Data[i] = nomatch
			} else {
				out.Missing[i] = true
			}
		}
		frame.Result = out
		return nil
	case "do_makeunique":
		text, err := characterData(x)
		if err != nil {
			return err
		}
		sep := "."
		if len(frame.Arguments) > 1 {
			if s, e := frameText(c, frame, 1); e == nil {
				sep = s
			}
		}
		seen := map[string]int{}
		out := &CharacterVector{Data: make([]string, len(text.Data)), Missing: append([]bool(nil), text.Missing...)}
		for i, s := range text.Data {
			candidate := s
			if seen[candidate] != 0 {
				// The first duplicate is .1, not .2.  Keep a separate
				// monotonic candidate counter so existing suffixed input is
				// skipped without changing the base occurrence count.
				for suffix := 1; ; suffix++ {
					candidate = s + sep + fmt.Sprint(suffix)
					if seen[candidate] == 0 {
						break
					}
				}
			}
			seen[candidate] = 1
			out.Data[i] = candidate
		}
		frame.Result = out
		return nil
	case "do_vhash":
		out := &IntegerVector{Data: make([]int64, len(keys))}
		for i, key := range keys {
			var h uint64 = 1469598103934665603
			for _, b := range []byte(key) {
				h ^= uint64(b)
				h *= 1099511628211
			}
			out.Data[i] = int64(h & 0x7fffffffffffffff)
		}
		frame.Result = out
		return nil
	}
	return fmt.Errorf("unhandled matching coordinate %s", planCoordinate(frame.Plan))
}

func orderPermutation(v Value) []int {
	keys := vectorKeys(v)
	order := make([]int, len(keys))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(i, j int) bool { return keys[order[i]] < keys[order[j]] })
	return order
}

func kernelOrderFamily(c *Context, frame *LoweringFrame) error {
	x, err := frameValue(c, frame, 0)
	if err != nil {
		return err
	}
	if frame.Plan.CEntry == "do_sort" {
		decreasing := false
		if len(frame.Arguments) > 1 {
			value, err := frameValue(c, frame, 1)
			if err != nil {
				return err
			}
			decreasing = scalarLogical(value)
		}
		frame.Result, err = sortValue(x, decreasing, nil)
		return err
	}
	order := orderPermutation(x)
	keys := vectorKeys(x)
	switch frame.Plan.CEntry {
	case "do_psort":
		frame.Result = takePositions(x, order)
	case "do_order":
		out := &IntegerVector{Data: make([]int64, len(order))}
		for i, p := range order {
			out.Data[i] = int64(p + 1)
		}
		frame.Result = out
	case "do_isunsorted", "do_sorted_fpass":
		unsorted := false
		for i := 1; i < len(keys); i++ {
			if keys[i] < keys[i-1] {
				unsorted = true
				break
			}
		}
		frame.Result = &LogicalVector{Data: []Logical{logicalFromBool(unsorted)}}
	case "do_rank":
		out := &DoubleVector{Data: make([]float64, len(order))}
		for start := 0; start < len(order); {
			end := start + 1
			for end < len(order) && keys[order[start]] == keys[order[end]] {
				end++
			}
			rank := float64(start+1+end) / 2
			for j := start; j < end; j++ {
				out.Data[order[j]] = rank
			}
			start = end
		}
		frame.Result = out
	case "do_xtfrm":
		frame.Result = takePositions(x, order)
	default:
		return fmt.Errorf("unhandled ordering coordinate %s", planCoordinate(frame.Plan))
	}
	return nil
}

func logicalFromBool(v bool) Logical {
	if v {
		return True
	}
	return False
}
