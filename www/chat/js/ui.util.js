"use strict";
//

/**
 * Returns true if X is within a Z range of Y
 *
 * @param {Number} x
 * @param {Number} y
 * @param {Number} z
 * @returns {Boolean}
 */
export function numsNear(x, y, z) {
    return Math.abs(x - y) < z;
}

/**
 * @param {Element} ele an element.
 * @returns {Boolean} true if 'ele' is scrolled to within 5px of the bottom of its scroll.
 */
export function ele_atBottom(ele) {
    return numsNear(ele.scrollTop, ele.scrollHeight - ele.clientHeight, 5);
}
