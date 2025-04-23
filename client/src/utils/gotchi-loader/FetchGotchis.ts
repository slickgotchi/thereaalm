import { ethers } from "ethers";
import { CustomiseOptions, customiseSvg } from "./customize-svg";
import { Options } from "./options";

export interface Aavegotchi {
    id: number;
    name: string;
    modifiedNumericTraits: number[];
    withSetsRarityScore: number;
    svgSet: GotchiSVGSet;
}

export interface GotchiSVGSet {
    id: string;
    front: string;
    left: string;
    right: string;
    back: string;
    anims?: {
        normal: string,
        happy: string,
        sad: string,
        mad: string,
    }
}

export async function fetchAavegotchis(account: string): Promise<Aavegotchi[]> {
    const coreQuery = `
    query ($owner: String!) {
      aavegotchis(where: { owner: $owner }, first: 100) {
        id
        name
        modifiedNumericTraits
        withSetsRarityScore
      }
    }
  `;
    const coreResponse = await fetch(
        "https://subgraph.satsuma-prod.com/tWYl5n5y04oz/aavegotchi/aavegotchi-core-matic/api",
        {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                query: coreQuery,
                variables: { owner: account.toLowerCase() },
            }),
        }
    );
    if (!coreResponse.ok) {
        throw new Error(
            `Core subgraph request failed: ${coreResponse.status} ${coreResponse.statusText}`
        );
    }
    const coreData = await coreResponse.json();
    if (!coreData?.data?.aavegotchis) {
        throw new Error("No Aavegotchis found in core subgraph");
    }

    return coreData.data.aavegotchis.map((g: any) => ({
        id: Number(g.id),
        name: g.name,
        modifiedNumericTraits: g.modifiedNumericTraits.map(Number),
        withSetsRarityScore: g.withSetsRarityScore,
        svgs: {},
        // svgs: gotchisMap.get(g.id) || { front: '', left: '', right: '', back: '' },
    }));
}



// New function to fetch all SVGs for a single Gotchi ID
export async function fetchSingleGotchiSVGs(
    gotchiID: string,
    pixels: number = 64
): Promise<GotchiSVGSet> {
    const svgQuery = `
    query ($id: ID!) {
      aavegotchi(id: $id) {
        svg
        left
        right
        back
      }
    }
  `;
    const response = await fetch(
        "https://subgraph.satsuma-prod.com/tWYl5n5y04oz/aavegotchi/aavegotchi-svg-matic/api",
        {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                query: svgQuery,
                variables: { id: gotchiID },
            }),
        }
    );
    if (!response.ok) {
        throw new Error(
            `SVG subgraph request failed: ${response.status} ${response.statusText}`
        );
    }
    const data = await response.json();
    // console.log(`SVG fetch response for ${gotchiID}:`, JSON.stringify(data));
    if (data?.data?.aavegotchi) {
        const g = data.data.aavegotchi;

        return {
            id: gotchiID,
            front: serializeAndSize(customiseSvg(g.svg, Options.Normal), pixels),
            left: serializeAndSize(customiseSvg(g.left || g.svg, Options.Normal), pixels),
            right: serializeAndSize(customiseSvg(g.right || g.svg, Options.Normal), pixels),
            back: serializeAndSize(customiseSvg(g.back || g.svg, Options.Normal), pixels),
            anims: {
                normal: serializeAndSize(customiseSvg(g.svg, Options.Normal), pixels),
                happy: serializeAndSize(customiseSvg(g.svg, Options.Happy), pixels),
                sad: serializeAndSize(customiseSvg(g.svg, Options.Sad), pixels),
                mad: serializeAndSize(customiseSvg(g.svg, Options.Mad), pixels),
            }
        };

        // return {
        //     id: gotchiID,
        //     front: removeBackground(removeShadow(g.svg, pixels), pixels),
        //     left: removeBackground(removeShadow(g.left || g.svg, pixels), pixels),
        //     right: removeBackground(removeShadow(g.right || g.svg, pixels), pixels),
        //     back: removeBackground(removeShadow(g.back || g.svg, pixels), pixels),
        // };
    } else {
        console.error(`No SVGs found for Gotchi ID ${gotchiID}`);
        return { id: gotchiID, front: "", left: "", right: "", back: "" }; // Fallback to empty strings—handled by placeholder in GameScene
    }
}

export async function fetchBulkGotchiSVGs(
    gotchiIDs: string[],
    pixels: number = 64
): Promise<GotchiSVGSet[]> {
    // Helper function to fetch in batches
    const fetchBatch = async (ids: string[]) => {
        const svgQuery = `
            query ($ids: [ID!]!) {
                aavegotchis(where: { id_in: $ids }) {
                    id
                    svg
                    left
                    right
                    back
                }
            }
        `;

        const response = await fetch(
            "https://subgraph.satsuma-prod.com/tWYl5n5y04oz/aavegotchi/aavegotchi-svg-matic/api",
            {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    query: svgQuery,
                    variables: { ids },
                }),
            }
        );
        if (!response.ok) {
            throw new Error(
                `Bulk SVG subgraph request failed: ${response.status}`
            );
        }
        const data = await response.json();
        if (!data?.data?.aavegotchis) {
            console.error("No SVG data found for the requested Gotchi IDs");
            return ids.map(() => ({
                id: "",
                svg: "",
                left: "",
                right: "",
                back: "",
            }));
        }

        return data.data.aavegotchis.map((g: any) => ({
            id: g.id,
            front: serializeAndSize(customiseSvg(g.svg, Options.Normal), pixels),
            left: serializeAndSize(customiseSvg(g.left || g.svg, Options.Normal), pixels),
            right: serializeAndSize(customiseSvg(g.right || g.svg, Options.Normal), pixels),
            back: serializeAndSize(customiseSvg(g.back || g.svg, Options.Normal), pixels),
            anims: {
                normal: serializeAndSize(customiseSvg(g.svg, Options.Normal), pixels),
                happy: serializeAndSize(customiseSvg(g.svg, Options.Happy), pixels),
                sad: serializeAndSize(customiseSvg(g.svg, Options.Sad), pixels),
                mad: serializeAndSize(customiseSvg(g.svg, Options.Mad), pixels),
            }
        }));

        // return data.data.aavegotchis.map((g: any) => ({
        //     id: g.id,
        //     svg: removeBackground(removeShadow(g.svg, pixels), pixels),
        //     left: removeBackground(removeShadow(g.left || g.svg, pixels), pixels),
        //     right: removeBackground(removeShadow(g.right || g.svg, pixels), pixels),
        //     back: removeBackground(removeShadow(g.back || g.svg, pixels), pixels),
        // }));
    };

    // Split gotchiIDs into batches of 100
    const batchSize = 100;
    const batches = [];
    for (let i = 0; i < gotchiIDs.length; i += batchSize) {
        batches.push(gotchiIDs.slice(i, i + batchSize));
    }

    // Fetch each batch and combine the results
    const allSVGs: GotchiSVGSet[] = [];
    for (const batch of batches) {
        const batchResult = await fetchBatch(batch);
        allSVGs.push(...batchResult);
    }

    // Return all fetched SVGs
    return allSVGs;
}

export const removeBackground = (
    svgString: string,
    pixels: number = 64
): string => {
    const parser = new DOMParser();
    const doc = parser.parseFromString(svgString, "image/svg+xml");
    const svg = doc.getElementsByTagName("svg")[0];

    // we nee to set these
    svg.setAttribute("width", `${pixels}`);
    svg.setAttribute("height", `${pixels}`);

    // Helper function to find the background group (checking multiple possible classes)
    const findBackgroundGroup = (element: Element): Element | null => {
        // Check all <g> elements, including nested ones
        const groups = element.getElementsByTagName("g");
        for (const group of groups) {
            // Look for common background classes or patterns
            if (
                group.classList.contains("gotchi-bg") || // Original class
                group.classList.contains("wearable-bg") || // New class from your SVG
                group.classList.contains("gotchi-wearable") || // Parent class
                // Optionally, check for a large fill area or specific color (e.g., purple background)
                (group.querySelector('path[fill="#9325ee"]') &&
                    group.querySelector('path[d="M0 0v64h64V0H0"]')) // Large rectangle covering the entire viewBox
            ) {
                return group;
            }
        }
        return null;
    };

    // Search for the background in the top-level SVG and its nested SVGs
    let background: Element | null = findBackgroundGroup(svg);
    if (!background) {
        // Check nested <svg> elements (e.g., inside <g> elements)
        const nestedSvgs = svg.getElementsByTagName("svg");
        for (const nestedSvg of nestedSvgs) {
            background = findBackgroundGroup(nestedSvg);
            if (background) break;
        }
    }

    if (background && svg.contains(background)) {
        // If the background is nested inside another SVG, we need to handle it carefully
        const parent = background.parentElement;
        if (parent && parent.tagName.toLowerCase() === "svg") {
            parent.removeChild(background);
        } else if (parent && parent.tagName.toLowerCase() === "g") {
            parent.removeChild(background);
        } else {
            svg.removeChild(background);
        }
    } else {
        console.log(
            "No background group found in SVG (checked gotchi-bg, wearable-bg, gotchi-wearable, or purple fill)"
        );
    }

    // Ensure the SVG is valid and serialize it
    const serializer = new XMLSerializer();
    return serializer.serializeToString(svg);
};

export const serializeAndSize = (svgString: string, pixels: number = 64): string => {
    const parser = new DOMParser();
    const doc = parser.parseFromString(svgString, "image/svg+xml");
    const svg = doc.getElementsByTagName("svg")[0];

    // Set dimensions
    svg.setAttribute("width", `${pixels}`);
    svg.setAttribute("height", `${pixels}`);

    // Serialize and return the modified SVG
    const serializer = new XMLSerializer();
    return serializer.serializeToString(svg);
}


export const removeShadow = (svgString: string, pixels: number = 64): string => {
    const parser = new DOMParser();
    const doc = parser.parseFromString(svgString, "image/svg+xml");
    const svg = doc.getElementsByTagName("svg")[0];

    // Set dimensions
    svg.setAttribute("width", `${pixels}`);
    svg.setAttribute("height", `${pixels}`);

    // Helper function to find the shadow group
    const findShadowGroup = (element: Element): Element | null => {
        // Check all <g> elements for the shadow class
        const groups = element.getElementsByTagName("g");
        for (const group of groups) {
            if (group.classList.contains("gotchi-shadow")) {
                return group;
            }
        }
        return null;
    };

    // Search for the shadow in the top-level SVG and its nested SVGs
    let shadow: Element | null = findShadowGroup(svg);
    if (!shadow) {
        // Check nested <svg> elements if shadow isn't found at top level
        const nestedSvgs = svg.getElementsByTagName("svg");
        for (const nestedSvg of nestedSvgs) {
            shadow = findShadowGroup(nestedSvg);
            if (shadow) break;
        }
    }

    // Remove the shadow if found
    if (shadow && svg.contains(shadow)) {
        const parent = shadow.parentElement;
        if (parent) {
            parent.removeChild(shadow);
        }
    } else {
        console.log("No shadow group found in SVG (checked for class 'gotchi-shadow')");
    }

    // Serialize and return the modified SVG
    const serializer = new XMLSerializer();
    return serializer.serializeToString(svg);
};