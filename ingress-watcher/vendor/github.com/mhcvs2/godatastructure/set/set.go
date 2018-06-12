package set

//interface
type Set interface {
	Add(e interface{}) bool
	Remove(e interface{})
	Clear()
	Contains(e interface{}) bool
	Len() int
	Same(other Set) bool
	Elements() []interface{}
	String() string
	Copy() Set
}

//Super function----------------------------------------------
//Determine whether one is a super set of other
func IsSuperset(one Set, other Set) bool {
	if one == nil || other == nil {
		return false
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0 || oneLen == otherLen {
		return false
	}
	if oneLen > 0 && otherLen == 0{
		return true
	}
	for _, v := range other.Elements() {
		if !one.Contains(v){
			return false
		}
	}
	return true
}

//Return union of two Set
func Union(one Set, other Set) Set {
	if other == nil && one == nil {
		return nil
	}
	if other == nil {
		return one.Copy()
	}
	if one == nil {
		return other.Copy()
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0{
		return other.Copy()
	}
	if otherLen == 0{
		return one.Copy()
	}
	copyset := other.Copy()
	for _, key := range one.Elements() {
		copyset.Add(key)
	}
	return copyset
}

//Return Intersect of two Set
func Intersect(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0 || otherLen == 0{
		return nil
	}
	if oneLen >= otherLen {
		copyset := other.Copy()
		for _, key := range copyset.Elements(){
			if !one.Contains(key){
				copyset.Remove(key)
			}
		}
		return copyset
	}else{
		copyset := one.Copy()
		for _, key := range copyset.Elements(){
			if !other.Contains(key){
				copyset.Remove(key)
			}
		}
		return copyset
	}
}

//Return Difference of Set compared to other
func Difference(one Set, other Set) Set {
	if other == nil && one == nil {
		return nil
	}
	if other == nil {
		return one.Copy()
	}
	if one == nil {
		return nil
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0 || otherLen == 0{
		return one.Copy()
	}
	copyset := one.Copy()
	for _, key := range copyset.Elements() {
		if other.Contains(key){
			copyset.Remove(key)
		}
	}
	return copyset
}

//Return Summetric Difference of two Set
func SummetricDifference(one Set, other Set) Set {
	if other == nil && one == nil {
		return nil
	}
	if other == nil {
		return one.Copy()
	}
	if one == nil {
		return other.Copy()
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0 {
		return other.Copy()
	}
	if otherLen == 0{
		return one.Copy()
	}
	union := Union(one, other)
	intersect := Intersect(one, other)
	return Difference(union, intersect)
}

func NewSimpleSet() Set {
	return NewHashSet()
}

func IsSet(value interface{}) bool {
	if _, ok := value.(Set); ok {
		return true
	}
	return false
}

