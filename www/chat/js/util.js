"use strict";
//

/**
 * @param {String} name
 * @param {String[][]} attrs
 * @param {Node[]} children
 */
export function create_element(name, attrs = [], children = [], events = []) {
    const ele = document.createElement(name);
    attrs.forEach((v) => { ele.setAttribute(v[0], v[1]); });
    children.forEach((v) => { ele.appendChild(v); });
    events.forEach((v) => { ele.addEventListener(v[0], v[1]); });
    return ele;
}

/**
 * @param {string} string
 * @returns {Text}
 */
export function dcTN(string) {
    return document.createTextNode(string);
}

/**
 * @param {number} x1
 * @param {number} x2
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
 * @param {number} x
 * @param {number} y
 * @param {number} z
 */
export function numsNear(x, y, z) {
    return Math.abs(x - y) < z;
}

/**
 * @param {Element} ele an element.
 * @returns {boolean} true if 'ele' is scrolled to within 5px of the bottom of its scroll.
 */
export function ele_atBottom(ele) {
    return numsNear(ele.scrollTop, ele.scrollHeight - ele.clientHeight, 5);
}

/**
 * @param {string} key
 * @param {string} value
 */
export function setDataBinding(key, value) {
    if (value === undefined || value === null) value = "";
    const e = document.querySelectorAll(`[data-bind="${key}"]`);
    if (e.length === 0) return;
    e.forEach((v) => { v.textContent = value; });
}

/**
 * @param {HTMLElement} el
 */
export function deActivateChild(el) {
    for (const item of el.children) {
        if (item.classList.contains("active")) {
            item.classList.remove("active");
        }
    }
}

/**
 * @param {HTMLElement} ele
 * @param {RegExp} regex
 * @param {Function} matcher function(string): Node
 */
export function safe_html_replace(ele, regex, matcher) {
    for (let i = 0; i < ele.childNodes.length; i++) {
        const item = ele.childNodes[i];
        if (item.nodeName !== "#text") {
            continue;
        }
        const fixed = item.textContent.split(regex).map((v) => {
            return regex.test(v) ? matcher(v) : dcTN(v);
        });
        if (fixed.length === 1) {
            continue;
        }
        for (const itn of fixed) {
            ele.insertBefore(itn, item);
        }
        item.remove();
        i += fixed.length - 1;
    }
}

/**
 * @param {HTMLDivElement} el
 * @param {number} ex
 */
export function popup_set_x(el, ex) {
    let x = ex + 24;
    const ew = el.offsetWidth;
    const ww = window.innerWidth;
    if (x + ew > ww) { x = ex - ew - 24; }
    el.style.left = `${x}px`;
}

/**
 * @param {HTMLDivElement} el
 * @param {number} ey
 */
export function popup_set_y(el, ey) {
    let y = ey - 24;
    const eh = el.offsetHeight;
    const wh = window.innerHeight;
    if (y + eh > wh) { y = wh - eh - 24; }
    el.style.top = `${y}px`;
}

/**
 * @param {HTMLElement} el
 * @param {string} field
 */
export function incr_data(el, field) {
    const v = parseInt(el.dataset[field], 10) | 0;
    el.dataset[field] = v + 1;
}

/**
 * @param {HTMLElement} el
 * @param {string} field
 */
export function decr_data(el, field) {
    const v = parseInt(el.dataset[field], 10) | 0;
    el.dataset[field] = v - 1;
}
