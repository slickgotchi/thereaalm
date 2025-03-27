interface HeapNode {
    value: number; // f-score
    data: any; // Node object
}

export class MinHeap {
    private heap: HeapNode[] = [];

    public insert(value: number, data: any) {
        this.heap.push({ value, data });
        this.bubbleUp(this.heap.length - 1);
    }

    public extractMin(): any {
        if (this.heap.length === 0) return null;
        if (this.heap.length === 1) return this.heap.pop()!.data;

        const min = this.heap[0].data;
        this.heap[0] = this.heap.pop()!;
        this.bubbleDown(0);
        return min;
    }

    public isEmpty(): boolean {
        return this.heap.length === 0;
    }

    private bubbleUp(index: number) {
        while (index > 0) {
            const parent = Math.floor((index - 1) / 2);
            if (this.heap[parent].value <= this.heap[index].value) break;
            [this.heap[parent], this.heap[index]] = [this.heap[index], this.heap[parent]];
            index = parent;
        }
    }

    private bubbleDown(index: number) {
        const length = this.heap.length;
        while (true) {
            let smallest = index;
            const left = 2 * index + 1;
            const right = 2 * index + 2;

            if (left < length && this.heap[left].value < this.heap[smallest].value) smallest = left;
            if (right < length && this.heap[right].value < this.heap[smallest].value) smallest = right;
            if (smallest === index) break;

            [this.heap[index], this.heap[smallest]] = [this.heap[smallest], this.heap[index]];
            index = smallest;
        }
    }

    // Add a method to check if a node exists (used in Pathfinder)
    public findNode(x: number, y: number): any {
        for (const node of this.heap) {
            if (node.data.x === x && node.data.y === y) return node.data;
        }
        return null;
    }
}