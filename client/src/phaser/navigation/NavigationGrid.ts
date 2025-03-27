export class NavigationGrid {
    private grid: boolean[][]; // true = passable, false = non-passable
    private width: number;
    private height: number;

    constructor(width: number, height: number) {
        this.width = width;
        this.height = height;
        this.grid = Array(height).fill(null).map(() => Array(width).fill(true)); // All passable by default
    }

    public isPassable(x: number, y: number): boolean {
        if (x < 0 || x >= this.width || y < 0 || y >= this.height) return false;
        return this.grid[y][x];
    }

    public setPassable(x: number, y: number, passable: boolean): void {
        if (x >= 0 && x < this.width && y >= 0 && y < this.height) {
            this.grid[y][x] = passable;
        }
    }

    public getWidth(): number {
        return this.width;
    }

    public getHeight(): number {
        return this.height;
    }
}