package world

import (
	"thereaalm/action"
	"thereaalm/action/buildingaction"
	"thereaalm/action/resourceaction"
	"thereaalm/entity"
	"thereaalm/types"
)

type GotchiBehaviorProfile func(g *entity.Gotchi) 

var MercenaryProfile GotchiBehaviorProfile = func(g *entity.Gotchi) {
	// mercenary actions
    g.AddActionToPlan(
        action.NewAttackAction(g, nil, 1, 
            &types.TargetSpec{
                TargetType: "lickquidator",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        action.NewAttackAction(g, nil, 1, 
            &types.TargetSpec{
                TargetType: "lickvoid",
                TargetCriterion: "nearest",
            }))

	// farmer actions
    g.AddActionToPlan(
        resourceaction.NewForageAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "fomoberrybush",
                TargetCriterion: "nearest",
            }))

	// minerjack actions
    g.AddActionToPlan(
        resourceaction.NewChopAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "kekwoodtree",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        resourceaction.NewMineAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "alphaslateboulders",
                TargetCriterion: "nearest",
            }))

	// builder actions
    g.AddActionToPlan(
        buildingaction.NewMaintainAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        buildingaction.NewRebuildAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    


	// general actions
    g.AddActionToPlan(
        action.NewSellAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "shop",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        action.NewRoamAction(g, nil, 0.1, nil))
}

var FarmerProfile GotchiBehaviorProfile = func(g *entity.Gotchi) {
	// mercenary actions
    g.AddActionToPlan(
        action.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickquidator",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        action.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickvoid",
                TargetCriterion: "nearest",
            }))
			
	// farmer actions
    g.AddActionToPlan(
        resourceaction.NewForageAction(g, nil, 1, 
            &types.TargetSpec{
                TargetType: "fomoberrybush",
                TargetCriterion: "nearest",
            }))

	// minerjack actions
    g.AddActionToPlan(
        resourceaction.NewChopAction(g, nil, 0.5, 
            &types.TargetSpec{
                TargetType: "kekwoodtree",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        resourceaction.NewMineAction(g, nil, 0.5, 
            &types.TargetSpec{
                TargetType: "alphaslateboulders",
                TargetCriterion: "nearest",
            }))

	// builder actions
    g.AddActionToPlan(
        buildingaction.NewMaintainAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        buildingaction.NewRebuildAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    


	// general actions
    g.AddActionToPlan(
        action.NewSellAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "shop",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        action.NewRoamAction(g, nil, 0.1, nil))
}

var MinerJackProfile GotchiBehaviorProfile = func(g *entity.Gotchi) {
	// mercenary actions
    g.AddActionToPlan(
        action.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickquidator",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        action.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickvoid",
                TargetCriterion: "nearest",
            }))
			
	// farmer actions
    g.AddActionToPlan(
        resourceaction.NewForageAction(g, nil, 0.5, 
            &types.TargetSpec{
                TargetType: "fomoberrybush",
                TargetCriterion: "nearest",
            }))

	// minerjack actions
    g.AddActionToPlan(
        resourceaction.NewChopAction(g, nil, 1, 
            &types.TargetSpec{
                TargetType: "kekwoodtree",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        resourceaction.NewMineAction(g, nil, 1, 
            &types.TargetSpec{
                TargetType: "alphaslateboulders",
                TargetCriterion: "nearest",
            }))

	// builder actions
    g.AddActionToPlan(
        buildingaction.NewMaintainAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        buildingaction.NewRebuildAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    


	// general actions
    g.AddActionToPlan(
        action.NewSellAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "shop",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        action.NewRoamAction(g, nil, 0.1, nil))
}

var BuilderProfile GotchiBehaviorProfile = func(g *entity.Gotchi) {
	// mercenary actions
    g.AddActionToPlan(
        action.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickquidator",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        action.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickvoid",
                TargetCriterion: "nearest",
            }))
			
	// farmer actions
    g.AddActionToPlan(
        resourceaction.NewForageAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "fomoberrybush",
                TargetCriterion: "nearest",
            }))

	// minerjack actions
    g.AddActionToPlan(
        resourceaction.NewChopAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "kekwoodtree",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        resourceaction.NewMineAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "alphaslateboulders",
                TargetCriterion: "nearest",
            }))

	// builder actions
    g.AddActionToPlan(
        buildingaction.NewMaintainAction(g, nil, 1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        buildingaction.NewRebuildAction(g, nil, 1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    


	// general actions
    g.AddActionToPlan(
        action.NewSellAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "shop",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        action.NewRoamAction(g, nil, 0.1, nil))
}

var ExplorerProfile GotchiBehaviorProfile = func(g *entity.Gotchi) {
	// mercenary actions
    g.AddActionToPlan(
        action.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickquidator",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        action.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickvoid",
                TargetCriterion: "nearest",
            }))
			
	// farmer actions
    g.AddActionToPlan(
        resourceaction.NewForageAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "fomoberrybush",
                TargetCriterion: "nearest",
            }))

	// minerjack actions
    g.AddActionToPlan(
        resourceaction.NewChopAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "kekwoodtree",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        resourceaction.NewMineAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "alphaslateboulders",
                TargetCriterion: "nearest",
            }))

	// builder actions
    g.AddActionToPlan(
        buildingaction.NewMaintainAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        buildingaction.NewRebuildAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    


	// general actions
    g.AddActionToPlan(
        action.NewSellAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "shop",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        action.NewRoamAction(g, nil, 1, nil))
}

