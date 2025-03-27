import { NavigationGrid } from "./NavigationGrid";
import { MinHeap } from "./MinHeap";
import { Direction } from "./TweenWorker";

interface Node {
    x: number;
    y: number;
    g: number;
    h: number;
    f: number;
    parent: Node | null;
}

export interface Waypoint {
    tileX: number;
    tileY: number;
    direction: "up" | "down" | "left" | "right" | "none";
}

export class Pathfinder {
    private grid: NavigationGrid;

    constructor(grid: NavigationGrid) {
        this.grid = grid;
    }

    public findPath(startX: number, startY: number, targetX: number, targetY: number): Waypoint[] {
        console.log(`[Pathfinder.findPath] Starting pathfinding from (${startX}, ${startY}) to (${targetX}, ${targetY})`);
        
        const openList = new MinHeap();
        const closedSet: Set<string> = new Set();
        const startNode: Node = { x: startX, y: startY, g: 0, h: this.heuristic(startX, startY, targetX, targetY), f: 0, parent: null };
        startNode.f = startNode.h;
        openList.insert(startNode.f, startNode);
        console.log(`[Pathfinder.findPath] Added start node (${startX}, ${startY}) to openList, f=${startNode.f}`);

        while (!openList.isEmpty()) {
            const current: Node = openList.extractMin();
            console.log(`[Pathfinder.findPath] Processing node (${current.x}, ${current.y}), f=${current.f}, g=${current.g}, h=${current.h}`);

            if (current.x === targetX && current.y === targetY) {
                console.log(`[Pathfinder.findPath] Target reached at (${current.x}, ${current.y})! Reconstructing path...`);
                return this.reconstructPath(current);
            }

            closedSet.add(`${current.x},${current.y}`);
            console.log(`[Pathfinder.findPath] Added (${current.x}, ${current.y}) to closedSet, size=${closedSet.size}`);

            const neighbors = this.getJumpPoints(current, targetX, targetY);
            console.log(`[Pathfinder.findPath] Found ${neighbors.length} jump points from (${current.x}, ${current.y}):`, neighbors);

            for (const neighbor of neighbors) {
                if (closedSet.has(`${neighbor.x},${neighbor.y}`)) {
                    console.log(`[Pathfinder.findPath] Skipping neighbor (${neighbor.x}, ${neighbor.y}) - already in closedSet`);
                    continue;
                }

                const tentativeG = current.g + this.manhattanDistance(current.x, current.y, neighbor.x, neighbor.y);
                const existing = openList.findNode(neighbor.x, neighbor.y);
                console.log(`[Pathfinder.findPath] Neighbor (${neighbor.x}, ${neighbor.y}), tentativeG=${tentativeG}, existing=${!!existing}`);

                if (!existing) {
                    const neighborNode: Node = { x: neighbor.x, y: neighbor.y, g: tentativeG, h: this.heuristic(neighbor.x, neighbor.y, targetX, targetY), f: 0, parent: current };
                    neighborNode.f = neighborNode.g + neighborNode.h;
                    openList.insert(neighborNode.f, neighborNode);
                    console.log(`[Pathfinder.findPath] Added new node (${neighborNode.x}, ${neighborNode.y}) to openList, f=${neighborNode.f}`);
                } else if (tentativeG < existing.g) {
                    console.log(`[Pathfinder.findPath] Updating existing node (${existing.x}, ${existing.y}) with better g: ${existing.g} -> ${tentativeG}`);
                    existing.g = tentativeG;
                    existing.f = existing.g + existing.h;
                    existing.parent = current;
                    openList.insert(existing.f, existing);
                }
            }
        }

        console.log(`[Pathfinder.findPath] Open list exhausted, no path found to (${targetX}, ${targetY})`);
        return [];
    }

    private heuristic(x1: number, y1: number, x2: number, y2: number): number {
        return Math.abs(x2 - x1) + Math.abs(y2 - y1);
    }

    private manhattanDistance(x1: number, y1: number, x2: number, y2: number): number {
        return Math.abs(x2 - x1) + Math.abs(y2 - y1);
    }

    private getJumpPoints(current: Node, targetX: number, targetY: number): Node[] {
        const directions = [];
        if (current.x < targetX) directions.push({ dx: 1, dy: 0, dir: "right" });
        if (current.x > targetX) directions.push({ dx: -1, dy: 0, dir: "left" });
        if (current.y < targetY) directions.push({ dx: 0, dy: 1, dir: "down" });
        if (current.y > targetY) directions.push({ dx: 0, dy: -1, dir: "up" });

        const jumpPoints: Node[] = [];

        console.log(`[Pathfinder.getJumpPoints] Checking jump points from (${current.x}, ${current.y}) toward (${targetX}, ${targetY})`);
        
        // Randomly shuffle directions to introduce variety
        for (let i = directions.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [directions[i], directions[j]] = [directions[j], directions[i]];
        }

        // Take one random step in a valid direction
        for (const { dx, dy, dir } of directions) {
            const nextX = current.x + dx;
            const nextY = current.y + dy;

            console.log(`[Pathfinder.getJumpPoints] Checking direction ${dir} to (${nextX}, ${nextY})`);
            if (this.grid.isPassable(nextX, nextY)) {
                const jumpPoint = { x: nextX, y: nextY, g: 0, h: 0, f: 0, parent: null };
                jumpPoints.push(jumpPoint);
                console.log(`[Pathfinder.getJumpPoints] Direction ${dir}: Added jump point (${nextX}, ${nextY})`);
                break; // Only add one random neighbor per step
            }
        }

        return jumpPoints;
    }

    private reconstructPath(endNode: Node): Waypoint[] {
        const waypoints: Waypoint[] = [];
        let current: Node | null = endNode;

        console.log(`[Pathfinder.reconstructPath] Starting path reconstruction from (${endNode.x}, ${endNode.y})`);
        while (current && current.parent) {
            const parent: any = current.parent;
            let direction: Direction = "none";
            if (current.x > parent.x) direction = "right";
            else if (current.x < parent.x) direction = "left";
            else if (current.y > parent.y) direction = "down";
            else if (current.y < parent.y) direction = "up";
            waypoints.unshift({ tileX: current.x, tileY: current.y, direction });
            console.log(`[Pathfinder.reconstructPath] Added waypoint (${current.x}, ${current.y}, ${direction})`);
            current = parent;
        }

        console.log(`[Pathfinder.reconstructPath] Path reconstructed with ${waypoints.length} waypoints:`, waypoints);
        return waypoints;
    }
}