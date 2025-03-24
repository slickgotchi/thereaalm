import { BaseEntity, EntitySnapshot } from "./BaseEntity";
import { EntityFactory } from "./EntityFactory";

export class EntityManager {
    private entities: Map<string, BaseEntity> = new Map();
    private scene: Phaser.Scene;

    constructor(scene: Phaser.Scene) {
        this.scene = scene;
    }

    updateEntities(entitySnapshots: EntitySnapshot[]) {
        const existingIds = new Set(this.entities.keys());

        entitySnapshots.forEach(snapshot => {
            if (this.entities.has(snapshot.id)) {
                this.entities.get(snapshot.id)!.update(snapshot);
                existingIds.delete(snapshot.id);
            } else {
                const entity = EntityFactory.create(this.scene, snapshot);
                this.entities.set(snapshot.id, entity);
            }
        });

        // Remove entities that no longer exist
        existingIds.forEach(id => {
            this.entities.get(id)!.destroy();
            this.entities.delete(id);
        });
    }
}
