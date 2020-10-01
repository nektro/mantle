"use strict";

//
import { create_element, dcTN, popup_set_x, popup_set_y } from "./../../../js/util.js";
import { el_xucm } from "./../../../js/ui.util.js";
import * as api from "./../../../js/api/index.js";
import * as ui from "./../../../js/ui.js";

/**
 * @param {HTMLUListElement} el
 * @param {number} ey
 */
function set_y(el, ey) {
    const y = ey - 24;
    const eh = el.offsetHeight;
    const wh = window.innerHeight;
    console.log(y,eh,wh);
    if (y + eh > wh) {
        el.style.bottom = "0";
        el.style.top = "initial";
    }
    if (eh > wh) {
        el.style.top = `-${ey-48}px`;
        el.style.bottom = "initial";
        el.style.maxHeight = `${wh-48}px`;
    }
}

//
customElements.define("x-user-contextmenu", class extends HTMLElement {
    constructor() {
        super();
        this.triggers = [];
        this.triggers.push(this.localName);
    }
    /**
     * @param {string} uid
     * @param {MouseEvent} e
     */
    async openWith(uid, e) {
        e.preventDefault();
        const user = await api.M.users.get(uid);
        const rolesU = (await user.getRoles()).map((v) => v.uuid);
        const rolesA = [...api.C.roles.values()].filter((v) => v.id !== undefined);

        this.removeAllChildren();
        this.appendChild(create_element("ul", [], [
            create_element("li", [], [create_element("a", [["href",`./../~${user.uuid}`],["target","_hello"]], [dcTN("Profile")]),]),
            create_element("li", [["data-requires","manage_roles"]], [
                create_element("a", [["class","more"]], [dcTN("Roles")]),
                create_element("ul", [], [
                    ...rolesA.map((v) => {
                        const ch = rolesU.includes(v.uuid) ? ["checked",""] : [];
                        return create_element("li", [], [create_element("a", [], [
                            create_element("label", [["for",`xucm-r-${v.uuid}`]], [dcTN(v.name)]),
                            create_element("input", [["id",`xucm-r-${v.uuid}`],["type","checkbox"],ch], [], [["change", async (ev) => {
                                const t = ev.target;
                                const r = t.id.split("-")[2];
                                if (ui.volatile.me.perms.manage_roles) {
                                    if (t.checked) {
                                        await api.M.users.update(uid, "add_role", r);
                                    } else {
                                        await api.M.users.update(uid, "remove_role", r);
                                    }
                                }
                            }]]),
                        ])]);
                    }),
                ]),
            ]),
            create_element("li", [["data-requires","manage_bans"]], [create_element("hr")]),
            create_element("li", [["data-requires","manage_bans"]], [create_element("a", [["class","danger"]], [dcTN("Kick")]),]),
            create_element("li", [["data-requires","manage_bans"]], [create_element("a", [["class","danger"]], [dcTN("Ban")]),]),
        ]));
        ui.refresh_permissions();
        popup_set_x(this, e.x);
        popup_set_y(this, e.y);
        for (const parent of document.querySelectorAll("ul li a")) {
            parent.addEventListener("mouseover", (_) => {
                if (_.target.localName !== "a") { return; }
                if (_.target.parentElement.children.length === 1) { return; }
                const el = _.target.nextElementSibling;
                set_y(el, e.y);
            });
        }
    }
});

//
document.addEventListener("click", (e) => {
    const p = e.target.path();
    for (const item of el_xucm.triggers) {
        for (const jtem of p) {
            if (jtem.matches(item)) {
                return;
            }
        }
    }
    el_xucm.removeAllChildren();
});
