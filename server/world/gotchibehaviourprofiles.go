package world

import (
	"thereaalm/action/buildingactions"
	"thereaalm/action/combatactions"
	"thereaalm/action/explorationactions"
	"thereaalm/action/resourceactions"
	"thereaalm/action/tradeactions"
	"thereaalm/entity"
	"thereaalm/types"
)

type GotchiBehaviorProfile func(g *entity.Gotchi) 

var MercenaryProfile GotchiBehaviorProfile = func(g *entity.Gotchi) {
	// mercenary actions
    g.AddActionToPlan(
        combatactions.NewAttackAction(g, nil, 1, 
            &types.TargetSpec{
                TargetType: "lickquidator",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        combatactions.NewAttackAction(g, nil, 1, 
            &types.TargetSpec{
                TargetType: "lickvoid",
                TargetCriterion: "nearest",
            }))

	// farmer actions
    g.AddActionToPlan(
        resourceactions.NewForageAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "fomoberrybush",
                TargetCriterion: "nearest",
            }))

	// minerjack actions
    g.AddActionToPlan(
        resourceactions.NewChopAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "kekwoodtree",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        resourceactions.NewMineAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "alphaslateboulders",
                TargetCriterion: "nearest",
            }))

	// builder actions
    g.AddActionToPlan(
        buildingactions.NewMaintainAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        buildingactions.NewRebuildAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    


	// general actions
    g.AddActionToPlan(
        tradeactions.NewSellAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "shop",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        explorationactions.NewRoamAction(g, nil, 0.1, nil))
}

var FarmerProfile GotchiBehaviorProfile = func(g *entity.Gotchi) {
	// mercenary actions
    g.AddActionToPlan(
        combatactions.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickquidator",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        combatactions.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickvoid",
                TargetCriterion: "nearest",
            }))
			
	// farmer actions
    g.AddActionToPlan(
        resourceactions.NewForageAction(g, nil, 1, 
            &types.TargetSpec{
                TargetType: "fomoberrybush",
                TargetCriterion: "nearest",
            }))

	// minerjack actions
    g.AddActionToPlan(
        resourceactions.NewChopAction(g, nil, 0.5, 
            &types.TargetSpec{
                TargetType: "kekwoodtree",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        resourceactions.NewMineAction(g, nil, 0.5, 
            &types.TargetSpec{
                TargetType: "alphaslateboulders",
                TargetCriterion: "nearest",
            }))

	// builder actions
    g.AddActionToPlan(
        buildingactions.NewMaintainAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        buildingactions.NewRebuildAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    


	// general actions
    g.AddActionToPlan(
        tradeactions.NewSellAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "shop",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        explorationactions.NewRoamAction(g, nil, 0.1, nil))
}

var MinerJackProfile GotchiBehaviorProfile = func(g *entity.Gotchi) {
	// mercenary actions
    g.AddActionToPlan(
        combatactions.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickquidator",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        combatactions.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickvoid",
                TargetCriterion: "nearest",
            }))
			
	// farmer actions
    g.AddActionToPlan(
        resourceactions.NewForageAction(g, nil, 0.5, 
            &types.TargetSpec{
                TargetType: "fomoberrybush",
                TargetCriterion: "nearest",
            }))

	// minerjack actions
    g.AddActionToPlan(
        resourceactions.NewChopAction(g, nil, 1, 
            &types.TargetSpec{
                TargetType: "kekwoodtree",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        resourceactions.NewMineAction(g, nil, 1, 
            &types.TargetSpec{
                TargetType: "alphaslateboulders",
                TargetCriterion: "nearest",
            }))

	// builder actions
    g.AddActionToPlan(
        buildingactions.NewMaintainAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        buildingactions.NewRebuildAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    


	// general actions
    g.AddActionToPlan(
        tradeactions.NewSellAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "shop",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        explorationactions.NewRoamAction(g, nil, 0.1, nil))
}

var BuilderProfile GotchiBehaviorProfile = func(g *entity.Gotchi) {
	// mercenary actions
    g.AddActionToPlan(
        combatactions.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickquidator",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        combatactions.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickvoid",
                TargetCriterion: "nearest",
            }))
			
	// farmer actions
    g.AddActionToPlan(
        resourceactions.NewForageAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "fomoberrybush",
                TargetCriterion: "nearest",
            }))

	// minerjack actions
    g.AddActionToPlan(
        resourceactions.NewChopAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "kekwoodtree",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        resourceactions.NewMineAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "alphaslateboulders",
                TargetCriterion: "nearest",
            }))

	// builder actions
    g.AddActionToPlan(
        buildingactions.NewMaintainAction(g, nil, 1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        buildingactions.NewRebuildAction(g, nil, 1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    


	// general actions
    g.AddActionToPlan(
        tradeactions.NewSellAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "shop",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        explorationactions.NewRoamAction(g, nil, 0.1, nil))
}

var ExplorerProfile GotchiBehaviorProfile = func(g *entity.Gotchi) {
	// mercenary actions
    g.AddActionToPlan(
        combatactions.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickquidator",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        combatactions.NewAttackAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "lickvoid",
                TargetCriterion: "nearest",
            }))
			
	// farmer actions
    g.AddActionToPlan(
        resourceactions.NewForageAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "fomoberrybush",
                TargetCriterion: "nearest",
            }))

	// minerjack actions
    g.AddActionToPlan(
        resourceactions.NewChopAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "kekwoodtree",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        resourceactions.NewMineAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "alphaslateboulders",
                TargetCriterion: "nearest",
            }))

	// builder actions
    g.AddActionToPlan(
        buildingactions.NewMaintainAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        buildingactions.NewRebuildAction(g, nil, 0.1,
            &types.TargetSpec{
                TargetType: "altar",
                TargetCriterion: "nearest",
            }))
    


	// general actions
    g.AddActionToPlan(
        tradeactions.NewSellAction(g, nil, 0.1, 
            &types.TargetSpec{
                TargetType: "shop",
                TargetCriterion: "nearest",
            }))
    g.AddActionToPlan(
        explorationactions.NewRoamAction(g, nil, 1, nil))
}

