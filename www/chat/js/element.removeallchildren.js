"use strict";
//
if (!('removeAllChildren' in Element.prototype)) {
    Element.prototype.removeAllChildren = function removeAllChildren() {
        if (this.children.length === 0) {
            return;
        }
        this.children[0].remove();
        this.removeAllChildren();
    };
}
