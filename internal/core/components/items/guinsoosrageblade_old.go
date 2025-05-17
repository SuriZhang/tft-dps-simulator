package items

// // GuinsoosRagebladeEffect_old tracks the stacking attack speed bonus from Guinsoo's Rageblade.
// type GuinsoosRagebladeEffect_old struct {
//     AttackSpeedPerStack float64 // The % AS gained per stack (e.g., 0.05 for 5%)
//     CurrentStacks       int     // Current number of stacks
//     // MaxStacks           int     // Optional: If there's a cap in the future
//     CurrentBonusAS      float64 // The total bonus AS currently provided by stacks
// }

// // GuinsoosRagebladeEffect_old creates a new GuinsoosRagebladeEffect_old component.
// func GuinsoosRagebladeEffect_old(asPerStack float64 /*, maxStacks int*/) *GuinsoosRagebladeEffect_old {
//     return &GuinsoosRagebladeEffect_old{
//         AttackSpeedPerStack: asPerStack,
//         CurrentStacks:       0,
//         // MaxStacks:           maxStacks,
//         CurrentBonusAS:      0.0,
//     }
// }

// // AddStack increments the stack count and updates the bonus AS.
// func (e *GuinsoosRagebladeEffect_old) IncrementStacks() {
//     // if e.MaxStacks > 0 && e.CurrentStacks >= e.MaxStacks {
//     //     return // Already at max stacks
//     // }
//     e.CurrentStacks++
//     e.CurrentBonusAS = float64(e.CurrentStacks) * e.AttackSpeedPerStack
// }

// // GetCurrentBonusAS returns the total bonus attack speed from current stacks.
// func (e *GuinsoosRagebladeEffect_old) GetCurrentBonusAS() float64 {
//     return e.CurrentBonusAS
// }

// // GetCurrentStacks returns the current number of stacks.
// func (e *GuinsoosRagebladeEffect_old) GetCurrentStacks() int {
//     return e.CurrentStacks
// }

// // ResetEffects resets the stacks and bonus AS.
// func (e *GuinsoosRagebladeEffect_old) ResetEffects() {
//     e.CurrentStacks = 0
//     e.CurrentBonusAS = 0.0
// }