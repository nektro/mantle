"use strict";
//

//
if (!("removeAllChildren" in Element.prototype)) {
    Element.prototype.removeAllChildren = function() {
        while (this.children.length > 0) {
            this.children[0].remove();
        }
    };
}
if (!("indexOfMe" in Element.prototype)) {
    Element.prototype.indexOfMe = function() {
        return Array.from(this.parentElement.children).indexOf(this);
    };
}
if (!("path" in Element.prototype)) {
    Element.prototype.path = function() {
        const p = [];
        let r = this;
        while (r !== null) {
            p.push(r);
            r = r.parentElement;
        }
        return p;
    };
}
