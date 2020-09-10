// base class for drag/droppable mementos
class Draggable {
    // getDragImage() - returns el
    // getDragData() - returns map of mime, data.

    // dropability (often) depends on type
    getType() {
      return null;
    }
};
