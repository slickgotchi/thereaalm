interface ResourceProps {
    scene: Phaser.Scene;
    x: number;
    y: number;
    iconTexture: string; // Key for the resource icon texture
    resourceCount: number; // Number of resources left
    trackingSprite: Phaser.GameObjects.Sprite;
}

export class ResourceIcon {
    private scene: Phaser.Scene;
    private resourceCount: number;
    private icon: Phaser.GameObjects.Image;
    private countText: Phaser.GameObjects.Text;
    private trackingSprite: Phaser.GameObjects.Sprite;

    constructor(props: ResourceProps) {
        const { scene, x, y, iconTexture, resourceCount, trackingSprite } = props;

        this.scene = scene;
        this.resourceCount = resourceCount;
        this.trackingSprite = trackingSprite;

        // Create the resource icon
        this.icon = this.scene.add.image(x, y, iconTexture);
        this.icon.setOrigin(0, 0);

        // Create the count text
        this.countText = this.scene.add.text(x + 16, y, `${resourceCount}`, {
            fontFamily: "Pixelar",
            fontSize: '16px',
            color: '#ffffff',
            stroke: '#000000',
            strokeThickness: 2
        });

        // Set depth to ensure visibility
        this.icon.setDepth(5000);
        this.countText.setDepth(5001);

        this.setPosition(x, y);
    }

    // Update the resource count and display
    public updateResourceCount(newCount: number): void {
        if (!this.trackingSprite) {
            this.destroy();
            return;
        }

        this.resourceCount = Math.max(0, newCount);
        this.countText.setText(`${this.resourceCount}`);

        this.setPosition(this.trackingSprite.x, this.trackingSprite.y);
    }

    // Update the position to follow the tracked sprite
    public setPosition(x: number, y: number): void {
        const offsetX = 1; // Offset to position icon above the sprite
        const offsetY = 1;

        this.icon.setPosition(x + offsetX, y + offsetY);
        this.countText.setPosition(x + offsetX + 16, y + offsetY - 8);
    }

    // Destroy the icon and text when no longer needed
    public destroy(): void {
        this.icon.destroy();
        this.countText.destroy();
    }
} 

// Example Usage:
// new ResourceIcon({ scene, x: 100, y: 200, iconTexture: 'wood_icon', resourceCount: 5, trackingSprite: someSprite });
