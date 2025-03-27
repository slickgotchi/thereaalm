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

    public findPath(startX: number, startY: number, targetX: number, targetY: number, useComplex: boolean): Waypoint[] {
        if (useComplex) {
            return this.findPathComplex(startX, startY, targetX, targetY);
        } else {
            return this.findPathSimple(startX, startY, targetX, targetY);
        }
    }

    private findPathComplex(startX: number, startY: number, targetX: number, targetY: number): Waypoint[] {
        const openList = new MinHeap();
        const closedSet: Set<string> = new Set();
        const startNode: Node = { x: startX, y: startY, g: 0, h: this.heuristic(startX, startY, targetX, targetY), f: 0, parent: null };
        startNode.f = startNode.h;
        openList.insert(startNode.f, startNode);

        while (!openList.isEmpty()) {
            const current: Node = openList.extractMin();

            if (current.x === targetX && current.y === targetY) {
                return this.reconstructPath(current);
            }

            closedSet.add(`${current.x},${current.y}`);

            const neighbors = this.getJumpPointsComplex(current, targetX, targetY);

            for (const neighbor of neighbors) {
                if (closedSet.has(`${neighbor.x},${neighbor.y}`)) continue;

                const tentativeG = current.g + this.manhattanDistance(current.x, current.y, neighbor.x, neighbor.y);
                const existing = openList.findNode(neighbor.x, neighbor.y);

                if (!existing) {
                    const neighborNode: Node = { x: neighbor.x, y: neighbor.y, g: tentativeG, h: this.heuristic(neighbor.x, neighbor.y, targetX, targetY), f: 0, parent: current };
                    neighborNode.f = neighborNode.g + neighborNode.h;
                    openList.insert(neighborNode.f, neighborNode);
                } else if (tentativeG < existing.g) {
                    existing.g = tentativeG;
                    existing.f = existing.g + existing.h;
                    existing.parent = current;
                    openList.insert(existing.f, existing);
                }
            }
        }

        return [];
    }

    private findPathSimple(startX: number, startY: number, targetX: number, targetY: number): Waypoint[] {
        const waypoints: Waypoint[] = [];
        let currentX = startX;
        let currentY = startY;

        while (currentX !== targetX || currentY !== targetY) {
            const directions = [];
            if (currentX < targetX) directions.push({ dx: 1, dy: 0, dir: "right" as Direction });
            if (currentX > targetX) directions.push({ dx: -1, dy: 0, dir: "left" as Direction });
            if (currentY < targetY) directions.push({ dx: 0, dy: 1, dir: "down" as Direction });
            if (currentY > targetY) directions.push({ dx: 0, dy: -1, dir: "up" as Direction });

            // Shuffle directions for randomness, but only valid moves toward target
            for (let i = directions.length - 1; i > 0; i--) {
                const j = Math.floor(Math.random() * (i + 1));
                [directions[i], directions[j]] = [directions[j], directions[i]];
            }

            // Pick the first direction (randomized, but always toward target)
            const { dx, dy, dir } = directions[0] || { dx: 0, dy: 0, dir: "none" as Direction };
            currentX += dx;
            currentY += dy;

            waypoints.push({ tileX: currentX, tileY: currentY, direction: dir });
        }

        return waypoints;
    }

    private heuristic(x1: number, y1: number, x2: number, y2: number): number {
        return Math.abs(x2 - x1) + Math.abs(y2 - y1);
    }

    private manhattanDistance(x1: number, y1: number, x2: number, y2: number): number {
        return Math.abs(x2 - x1) + Math.abs(y2 - y1);
    }

    private getJumpPointsComplex(current: Node, targetX: number, targetY: number): Node[] {
        const directions = [
            { dx: 1, dy: 0, dir: "right" },
            { dx: -1, dy: 0, dir: "left" },
            { dx: 0, dy: 1, dir: "down" },
            { dx: 0, dy: -1, dir: "up" },
        ];
        const jumpPoints: Node[] = [];

        for (let i = directions.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [directions[i], directions[j]] = [directions[j], directions[i]];
        }

        for (const { dx, dy } of directions) {
            const nextX = current.x + dx;
            const nextY = current.y + dy;

            if (this.grid.isPassable(nextX, nextY)) {
                const jumpPoint = { x: nextX, y: nextY, g: 0, h: 0, f: 0, parent: null };
                jumpPoints.push(jumpPoint);
            }
        }

        return jumpPoints;
    }

    private reconstructPath(endNode: Node): Waypoint[] {
        const waypoints: Waypoint[] = [];
        let current: Node | null = endNode;

        while (current && current.parent) {
            const parent: any = current.parent;
            let direction: Direction = "none";
            if (current.x > parent.x) direction = "right";
            else if (current.x < parent.x) direction = "left";
            else if (current.y > parent.y) direction = "down";
            else if (current.y < parent.y) direction = "up";
            waypoints.unshift({ tileX: current.x, tileY: current.y, direction });
            current = parent;
        }

        return waypoints;
    }
}