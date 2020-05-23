"use strict";

//
import { create_element, dcTN } from "./../util.js";
// import { output } from "./../ui.util.js";
import * as ui from "./../ui.js";
import * as api from "./../api/index.js";

//
customElements.define("x-user-dialog", class extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        //
    }
    /**
     * @param {string} uid
     * @param {MouseEvent} e
     */
    async openWith(uid, e) {
        const userN = await api.M.users.get(uid);
        this.removeAllChildren();
        this.appendChild(create_element("div", null, [
            create_element("h2", null, [dcTN(userN.nickname)]),
            create_element("h2", null, [
                create_element("span", null, [dcTN(userN.name)]),
                dcTN("#"),
                create_element("span", null, [dcTN(userN.id)]),
            ]),
            create_element("div", [["id","pp_uuid"]], [dcTN(userN.uuid)]),
            create_element("div", null, [
                dcTN("Provider: "),
                create_element("span", null, [dcTN(userN.provider)]),
            ]),
            create_element("hr"),
            create_element("div", null, [dcTN("Roles")]),
            create_element("ol"),
            create_element("div", [["data-requires","manage_roles"]], [
                create_element("i", [["class","plus icon"]]),
                create_element("ol"),
            ]),
        ]));
        for (const item of api.C.roles.values()) {
            if (item.id === undefined) { continue; }
            //
            const nEl2 = create_element("li", [["data-role",item.uuid],["class","bg-bf"]], [dcTN(item.name)]);
            nEl2.addEventListener("click", (e) => {
                if (!ui.volatile.me.perms.manage_roles) return;
                const et = e.target;
                const rid = et.dataset.role;
                const uid = document.querySelector("#pp_uuid").textContent;
                this.toggleRole(rid);
                return api.M.users.update(uid,"remove_role",rid);
            });
            this.children[0].querySelectorAll("ol")[0].appendChild(nEl2);
            //
            const nEl3 = create_element("li", [["data-role",item.uuid]], [dcTN(item.name)]);
            nEl3.addEventListener("click", (e) => {
                if (!ui.volatile.me.perms.manage_roles) return;
                const et = e.target;
                const rid = et.dataset.role;
                const uid = this.querySelector("#pp_uuid").textContent;
                this.toggleRole(rid);
                return api.M.users.update(uid,"add_role",rid);
            });
            this.children[0].querySelectorAll("ol")[1].appendChild(nEl3);
        }
        for (const item of await userN.getRoles()) {
            if (item.id === undefined) { continue; }
            this.toggleRole(item.uuid);
        }
        set_x(this.children[0], e.x);
        set_y(this.children[0], e.y);
    }
    toggleRole(uid) {
        for (const item of this.querySelectorAll("ol")) {
            item.querySelector(`[data-role="${uid}"]`).classList.toggle("active");
        }
    }
});

/**
 * @param {HTMLDivElement} el
 * @param {number} ex
 */
function set_x(el, ex) {
    let x = ex + 24;
    // const ew = el.offsetWidth;
    // const ww = window.inneroffsetWidth;
    // if (x + ew > ww) { x = ww - ew - 24; }
    el.style.left = `${x}px`;
}

/**
 * @param {HTMLDivElement} el
 * @param {number} ey
 */
function set_y(el, ey) {
    let y = ey - 24;
    const eh = el.offsetHeight;
    const wh = window.innerHeight;
    if (y + eh > wh) { y = wh - eh - 24; }
    el.style.top = `${y}px`;
}

//
document.addEventListener("click", (e) => {
    const p = e.target.path();
    const ud = document.querySelector("x-user-dialog");
    if (p.includes(ud)) { return; }
    if (p[0].classList.contains("usr")) { return; }
    ud.removeAllChildren();
});
