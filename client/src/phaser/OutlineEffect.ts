interface OutlineEffectConfig {
    scene: Phaser.Scene;
    target: Phaser.GameObjects.Sprite;
    color?: number;
    offset?: number;
}

export class OutlineEffect {
    private scene: Phaser.Scene;
    private target: Phaser.GameObjects.Sprite;
    private outlines: Phaser.GameObjects.Sprite[] = [];
    private color: number;
    private offset: number;
    private isDestroyed: boolean = false;

    constructor({ scene, target, color = 0xffffff, offset = 2 }: OutlineEffectConfig) {
        this.scene = scene;
        this.target = target;
        this.color = color;
        this.offset = offset;

        this.createOutlineSprites();
    }

    private createOutlineSprites(): void {
        this.destroyOutlinesOnly();

        const positions = [
            { x: -this.offset, y: 0 },  // Left
            { x: this.offset, y: 0 },   // Right
            { x: 0, y: -this.offset },  // Up
            { x: 0, y: this.offset }    // Down
        ];

        for (const pos of positions) {
            const outline = this.scene.add.sprite(this.target.x, this.target.y, this.target.texture.key)
                .setTintFill(this.color)
                .setDepth(this.target.depth - 1)
                .setVisible(false)
                .setOrigin(0, 0);

            this.outlines.push(outline);
        }
    }

    public rebuild(): void {
        if (this.isDestroyed) return;
        this.createOutlineSprites();
    }

    public showOutline(): void {
        if (this.isDestroyed) return;

        this.outlines.forEach((outline, index) => {
            const offset = this.getOffset(index);
            outline.setPosition(this.target.x + offset.x, this.target.y + offset.y);
            outline.setVisible(true);
        });
    }

    public hideOutline(): void {
        if (this.isDestroyed) return;

        this.outlines.forEach(outline => outline.setVisible(false));
    }

    private getOffset(index: number): { x: number; y: number } {
        const offsets = [
            { x: -this.offset, y: 0 },
            { x: this.offset, y: 0 },
            { x: 0, y: -this.offset },
            { x: 0, y: this.offset }
        ];
        return offsets[index];
    }

    private destroyOutlinesOnly(): void {
        this.outlines.forEach(outline => outline.destroy());
        this.outlines = [];
    }

    public destroy(): void {
        if (this.isDestroyed) return;

        this.destroyOutlinesOnly();

        this.scene = null!;
        this.target = null!;
        this.isDestroyed = true;
    }
}
