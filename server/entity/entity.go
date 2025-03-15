package entity

import "thereaalm/entity/component"

type Entity interface {
    Update()
    GetPosition() *component.Position // Returns the entity's position
    // Add more methods as needed (e.g., GetID(), GetType())
}