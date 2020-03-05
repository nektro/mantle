"use strict";
//

/**
 * @param {String} name
 * @param {String[][]} attrs
 * @param {Node[]} children
 * @returns {HTMLElement}
 */
export function create_element(name, attrs, children) {
    const ele = document.createElement(name);
    (attrs || []).forEach((v) => { ele.setAttribute(v[0], v[1]); });
    (children || []).forEach((v) => { ele.appendChild(v); });
    return ele;
}

/**
 * @param {String} string
 * @returns {Text}
 */
export function dcTN(string) {
    return document.createTextNode(string);
}

/**
 * @param {Number} x1
 * @param {Number} x2
 * @returns {Number[]}
 */
export function numsBetween(x1, x2) {
    if (x1 === x2) return [x1];
    const res = [];
    //
    if (x1 > x2) {
        for (let i = x2; i <= x1; i++) {
            res.push(i);
        }
    }
    if (x2 > x1) {
        for (let i = x1; i <= x2; i++) {
            res.push(i);
        }
    }
    return res;
}

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

/**
 * @param {String} key
 * @param {String} value
 */
export function setDataBinding(key, value) {
    const e = document.querySelector(`[data-bind="${key}"]`);
    if (e === null) return;
    e.textContent = value;
}
