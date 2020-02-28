"use strict";
//
if (!("removeAllChildren" in Element.prototype)) {
    Element.prototype.removeAllChildren = function() {
        while (this.children.length > 0) {
            this.children[0].remove();
        }
    };
}
if (!("isInViewport" in Element.prototype)) {
    Element.prototype.isInViewport = function() {
        const b = this.getBoundingClientRect();
        // console.log(b);
        return b.top > 0 && b.left > 0;
    };
}
