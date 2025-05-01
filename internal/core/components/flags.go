package components

// CanAbilityCritFromItems is a marker component indicating the entity's abilities can critically strike
// *due to an item effect* (like JG or IE).
type CanAbilityCritFromItems struct{}

// CanAbilityCritFromTraits is a marker component indicating the entity's abilities can critically strike
// *as part of their base kit*.
type CanAbilityCritFromTraits struct{}

// CanAbilityCritFromAugments is a marker component indicating the entity's abilities can critically strike
// *as part of their base kit*.
type CanAbilityCritFromAugments struct{}

// ImmuneToCC is a marker component indicating the entity is immune to crowd control effects. 
// TODO: add this for quicksilver
type ImmuneToCC struct{}