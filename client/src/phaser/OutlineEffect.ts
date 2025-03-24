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

    constructor({ scene, target, color = 0xffffff, offset = 2 }: OutlineEffectConfig) {
        this.scene = scene;
        this.target = target;
        this.color = color;
        this.offset = offset;

        this.createOutlineSprites();
        // this.setupInteraction();
    }

    private createOutlineSprites(): void {
        this.destroy();

        const positions = [
            { x: -this.offset, y: 0 },  // Left
            { x: this.offset, y: 0 },   // Right
            { x: 0, y: -this.offset },  // Up
            { x: 0, y: this.offset }    // Down
        ];

        for (const pos of positions) {
            const outline = this.scene.add.sprite(this.target.x, this.target.y, this.target.texture.key)
                .setTintFill(this.color)
                .setDepth(this.target.depth - 1) // Ensure it's behind the target
                .setVisible(false)
                .setOrigin(0,0);

            this.outlines.push(outline);
        }
    }

    public rebuild() {
        this.createOutlineSprites();
    }

    public showOutline(): void {
        this.outlines.forEach((outline, index) => {
            const offset = this.getOffset(index);
            outline.setPosition(this.target.x + offset.x, this.target.y + offset.y);
            outline.setVisible(true);
        });
    }

    public hideOutline(): void {
        this.outlines.forEach(outline => outline.setVisible(false));
    }

    private getOffset(index: number): { x: number; y: number } {
        const offsets = [
            { x: -this.offset, y: 0 },  // Left
            { x: this.offset, y: 0 },   // Right
            { x: 0, y: -this.offset },  // Up
            { x: 0, y: this.offset }    // Down
        ];
        return offsets[index];
    }

    public destroy(): void {
        this.outlines.forEach(outline => outline.destroy());
        this.outlines = [];
    }
} 

// Usage:
// const outline = new OutlineEffect({ scene: this, target: yourSprite });
