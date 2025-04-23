import { mouths, eyes, defaultGotchi } from "./svg-prefabs";
import { Tuple, AavegotchiObject } from "./types";


/**
 * Removes background from Aavegotchi SVG
 * @param {string} svg - SVG you want to customise
 * @returns {string} Returns customised SVG
 */
export const removeBG = (svg: string) => {
  const styledSvg = svg.replace("<style>", "<style>.gotchi-bg,.wearable-bg{display: none}");
  return styledSvg;
};

/**
 * Removes shadow from Aavegotchi SVG
 * @param {string} svg - SVG you want to customise
 * @returns {string} Returns customised SVG
 */
 export const removeShadow = (svg: string) => {
  const styledSvg = svg.replace("<style>", "<style>.gotchi-shadow{display: none}");
  return styledSvg;
};

export const scale = (svg: string, scale: {x: number, y: number}) => {
  let styledSvg = svg.replace(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">`,
  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
    <g transform="translate(32,32) scale(${scale.x},${scale.y}) translate(-32,-32)">
  `);

  styledSvg = styledSvg.slice(0,styledSvg.length-6) +"</g></svg>";

  return styledSvg;
}

export const skew = (svg: string, skew: {x: number, y: number}) => {
  let styledSvg = svg.replace(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">`,
  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
    <g transform="translate(32,32) skewX(${skew.x}) skewY(${skew.y}) translate(-32,-32)">
  `);

  styledSvg = styledSvg.slice(0,styledSvg.length-6) +"</g></svg>";

  return styledSvg;
}

export const rotate = (svg: string, angle: number) => {
  let styledSvg = svg.replace(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">`,
  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
    <g transform="rotate(${angle},32,32)">
  `);

  styledSvg = styledSvg.slice(0,styledSvg.length-6) +"</g></svg>";

  return styledSvg;
}

export const mirrorY = (svg: string) => {
  let styledSvg = svg.replace(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">`,
  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
    <g transform="scale(1,-1) translate(0,-64)">
  `);

  styledSvg = styledSvg.slice(0,styledSvg.length-6) +"</g></svg>";

  return styledSvg;
}

export const mirrorX = (svg: string) => {
  let styledSvg = svg.replace(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">`,
  `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
    <g transform="scale(-1,1) translate(-64,0)">
  `);

  styledSvg = styledSvg.slice(0,styledSvg.length-6) +"</g></svg>";

  return styledSvg;
}

/**
 * Adds Keyframe animation to SVG. (NOT TO BE USED IN IN GAME SPRITESHEET)
 * @param {string} svg - SVG you want to customise
 * @returns {string} Returns customised SVG
 */
export const bounceAnimation = (svg: string) => {
  const style = `
    @keyframes downHands {
      from {
        --hand_translateY: -1px;
      }
      to {
        --hand_translateY: -0.5px;
      }
    }
    @keyframes up {
      from {
        transform: translate(0px, 0);
      }
      to {
        transform: translate(0px, -1px);
      }
    }
    @keyframes down {
      from {
        transform: translate(0px, 0);
      }
      to {
        transform: translate(0px, 1px);
      }
    }
    svg {
      animation-name:down;
      animation-duration:1s;
      animation-iteration-count: infinite;
      animation-timing-function: linear;
      animation-timing-function: steps(1);
    }
    .gotchi-shadow {
      animation: up 1s infinite linear steps(2);
      animation-name:up;
      animation-duration:1s;
      animation-iteration-count: infinite;
      animation-timing-function: linear;
      animation-timing-function: steps(2);
    }
    .gotchi-wearable {
      animation-name:down;
      animation-duration:1s;
      animation-iteration-count: infinite;
      animation-timing-function: linear;
      animation-timing-function: steps(1);
    }
    .gotchi-handsDownClosed, .gotchi-handsUp, .gotchi-handsDownOpen, .gotchi-handsDownClosed, .gotchi-body, .gotchi-eyeColor, .gotchi-collateral, .gotchi-cheek, .gotchi-primary-mouth, .gotchi-wearable, .gotchi-sleeves   {
      animation-name:down;
      animation-duration:1s;
      animation-iteration-count: infinite;
      animation-timing-function: linear;
      animation-timing-function: steps(2);
    }
    .wearable-hand {
      animation-name:downHands !important;
      animation-duration:1s;
      animation-iteration-count: infinite;
      animation-timing-function: linear;
      animation-timing-function: steps(2);
      transform: translate(0, var(--hand_translateY));
    }
    .wearable-bg {
      animation-name: none;
    }
  `;
  const styledSvg = svg.replace("<style>", `<style>${style}`);
  return styledSvg;
};


/**
 * Adds SVG styling to Aavegotchi to raise its arms
 * @param {string} svg - SVG you want to customise
 * @param {{left?: number, right?: number}} arms - Wearable id of arms for unique animations 
 * @returns {string} Returns customised SVG
 */
export const raiseHands = (svg: string, arms?: {left?: number, right?: number}) => {
  const leftArm = (arms?.left && [207, 217, 223].includes(arms?.left)) ? `
      .wearable-hand-left {
        transform: translateY(calc(14px + var(--hand_translateY, -4px))) scaleY(-1);
        transform-origin: 50% 50%;
      }
    ` : ''
  const rightArm = (arms?.right && [207, 217, 223].includes(arms?.right)) ? `
    .wearable-hand-right {
      transform: translateY(calc(14px + var(--hand_translateY, -4px))) scaleY(-1);
      transform-origin: 50% 50%;
    }
  ` : ``

  const style = `
    .gotchi-handsDownClosed {
      display:none !important;
    }
    .gotchi-handsDownOpen {
      display:none !important;
    }
    .gotchi-handsUp {
      display:block !important;
    }
    .gotchi-sleeves {
      display:none !important;
    }
    .gotchi-sleeves-up {
      display:block !important;
    }
    .wearable-hand {
      transform: translateY(var(--hand_translateY, -4px));
    }
    ${leftArm}
    ${rightArm}
  `;

  const styledSvg = svg.replace("<style>", `<style>${style}`);
  return styledSvg;
};

/**
 * Adds SVG styling to Aavegotchi so it appears to float higher
 * @param {string} svg - SVG you want to customise
 * @returns {string} Returns customised SVG
 */
export const addIdleUp = (svg: string): string => {
  const styledSvg = svg.replace(
    "<style>",
    "<style>.gotchi-shadow {transform: translateY(1px);}.gotchi-wearable,.gotchi-handsDownClosed,.gotchi-handsUp,.gotchi-handsDownOpen,.gotchi-handsDownClosed,.gotchi-body,.gotchi-eyeColor,.gotchi-collateral,.gotchi-cheek,.gotchi-primary-mouth,.gotchi-wearable,.gotchi-sleeves {transform: translateY(-2px);}"
  );
  return styledSvg;
};

interface ReplaceEyes {
  target: "eyes";
  replaceSvg: keyof typeof eyes;
}

interface ReplaceMouth {
  target: "mouth";
  replaceSvg: keyof typeof mouths;
}

type ReplaceElement = ReplaceEyes | ReplaceMouth;

/**
 * Replaces a layer in the Aavegotchi SVG with custom SVG data
 * @param {string} svg - SVG you want to customise
 * @param {ReplaceElement} element - target of element you want to replace + element you want to replace it with
 * @returns {string} Returns customised SVG
 */
export function replaceParts(svg: string, element: ReplaceElement) {
  const doc = document.createDocumentFragment();
  const wrapper = document.createElement("svg");
  wrapper.setAttribute("xmlns", "http://www.w3.org/2000/svg");
  wrapper.setAttribute("viewbox", "0 0 64 64");
  wrapper.innerHTML = svg;
  doc.appendChild(wrapper);

  const targetClass =
    element.target === "eyes" ? "g.gotchi-eyeColor" : "g.gotchi-primary-mouth";
  const textnodes = doc.querySelectorAll(targetClass);

  textnodes.forEach(function (txt) {
    txt.innerHTML =
      element.target === "eyes"
        ? eyes[element.replaceSvg]
        : mouths[element.replaceSvg];
    //txt.parentNode?.replaceChild(el, txt);
  });
  const div = document.createElement("svg");
  div.appendChild(doc);
  return div.innerHTML;
}


export type CustomiseOptions = {
  removeBg?: boolean,
  eyes?: keyof typeof eyes,
  mouth?: keyof typeof mouths,
  float?: boolean,
  animate?: boolean,
  armsUp?: boolean,
  removeShadow?: boolean,
  scale?: {
    x: number,
    y: number
  },
  skew?: {
    x: number,
    y: number
  }
  mirrorY?: boolean,
  mirrorX?: boolean,
  rotate?: number,
}

/**
 * Customise Aavegotchi SVG
 * @param {string} svg - SVG you want to customise
 * @param {CustomiseOptions} options - Properties you want to change
 * @param {Tuple<number, 16>} equipped - Equipped wearables (Only necessary for raised mechanical arms)
 * @returns {string} Returns customised SVG
 */
export const customiseSvg = (svg: string, options: CustomiseOptions, equipped?: Tuple<number, 16>): string => {
  let styledSvg = svg;
  // remove bg and shadow as a default
  styledSvg = removeBG(styledSvg);
  styledSvg = removeShadow(styledSvg);

  // go through keys
  (Object.keys(options) as Array<keyof typeof options>).map((option) => {
    const value = options[option];
    if (value) {
      switch (option) {
        case 'removeBg':
          return styledSvg = removeBG(styledSvg);
        case 'eyes':
          return styledSvg = replaceParts(styledSvg, {target: option, replaceSvg: value as keyof typeof eyes});
        case 'mouth':
          return styledSvg = replaceParts(styledSvg, {target: option, replaceSvg: value as keyof typeof mouths});
        case 'animate':
          return styledSvg = bounceAnimation(styledSvg);
        case 'float':
          return styledSvg = addIdleUp(styledSvg);
        case 'armsUp':
          return styledSvg = raiseHands(styledSvg, equipped ? {left: equipped[4], right: equipped[5]} : undefined);
        case 'removeShadow':
          return styledSvg = removeShadow(styledSvg);
        case 'scale':
          return styledSvg = scale(styledSvg, options.scale ? {x: options.scale.x, y: options.scale.y} : {x: 1, y: 1});
        case 'skew':
          return styledSvg = skew(styledSvg, options.skew ? {x: options.skew.x, y: options.skew.y} : {x: 0, y: 0});
        case 'mirrorY':
          return styledSvg = mirrorY(styledSvg);
        case 'mirrorX':
          return styledSvg = mirrorX(styledSvg);
        case 'rotate':
          return styledSvg = rotate(styledSvg, options.rotate ? options.rotate : 0);
        default:
          return styledSvg;
      }
    }
  })
  return styledSvg;
}