import { ethers } from 'ethers';
export interface Tuple<T extends unknown, L extends number> extends Array<T> {
  0: T;
  length: L;
}

export interface Gotchi {
  tokenId: number;
  name: string;
  kinship: number;
  experience: number;
  level: number;
  nrg: number;
  agg: number;
  spk: number;
  brn: number;
  eys: number;
  eyc: number;
  svgs: string[];
  wearables: number[];
  blobUrl: string; 
}


export interface AavegotchiGameObject extends AavegotchiObject {
  spritesheetKey: string;
  svg: Tuple<string, 4>;
}

export interface AavegotchiObject extends AavegotchiContractObject {
  svg?: Tuple<string, 4>;
}

export interface AavegotchiContractObject {
  // Only in subgraph
  // withSetsNumericTraits: Tuple<number, 6>;
  // id: string;
  // withSetsRarityScore: number;
  // owner: {
  //   id: string;
  // };

  // collateral: string;
  name: string;
  modifiedNumericTraits: number[];
  level: number;
  // numericTraits: number[];
  // owner: string;
  // randomNumber: string;
  status: number;
  tokenId: ethers.BigNumberish;
  // items: ItemsAndBalances[];
  equippedWearables: Tuple<number, 16>;
  // experience: ethers.BigNumber;
  // hauntId: ethers.BigNumber;
  kinship: ethers.BigNumberish;
  experience: ethers.BigNumberish;
  // lastInteracted: string;
  // level: ethers.BigNumber;
  // toNextLevel: ethers.BigNumber;
  // stakedAmount: ethers.BigNumber;
  // minimumStake: ethers.BigNumber;
  // usedSkillPoints: ethers.BigNumber;
  // escrow: string;
  // baseRarityScore: ethers.BigNumber;
  // modifiedRarityScore: ethers.BigNumber;
  // locked: boolean;
  // unlockTime: string;
}

export interface ItemsAndBalances {
  itemType: ItemObject;
  itemId: ethers.BigNumberish;
  balance: ethers.BigNumberish;
}

export interface ItemObject {
  allowedCollaterals: number[];
  canBeTransferred: boolean;
  canPurchaseWithGhst: boolean;
  description?: string;
  category: number;
  experienceBonus: string;
  ghstPrice: ethers.BigNumberish;
  kinshipBonus: string;
  maxQuantity: ethers.BigNumberish;
  minLevel: string;
  name: string;
  rarityScoreModifier: string;
  setId: string;
  slotPositions: boolean[];
  svgId: number;
  totalQuantity: number;
  traitModifiers: number[];
}

export interface SubmitScoreReq {
  name: string,
  tokenId: string,
}

export interface HighScore {
  tokenId: string,
  score: number,
  name: string,
}

export interface CustomError extends Error {
  status?: number;
}

// custom types for my game
export interface LevelConfig {
  levelNumber: number,
  gridObjectLayout: Array<number[]>,
  levelHeader: string,
  levelDescription: string,
  gridColour: number,
  posBezPrev: number[],
  posBtn: number[],
  posBezNext: number[],
  // pos: number[],
  // curveThisPos: number[],
  // curvePrevPos: number[],
  actionsAvailable: number,
  maxPointsPossible: number,
}



export type SpriteMatrix = Array<Array<string>>;

export interface Spritesheet {
  src: string,
  dimensions: {
    width: number,
    height: number,
    x: number,
    y: number,
  }
}