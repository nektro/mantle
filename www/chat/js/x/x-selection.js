"use strict";
//
import { create_element, dcTN, deActivateChild } from "./../util.js";
import * as api from "./../api/index.js";

//
customElements.define("x-selection", class extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        if (!this.children[0].classList.contains("static")) {
            $(this.children[0]).sortable({
                // jshint -W098
                stop: (ev, ue) => {
                    const c = this.parentElement.getAttribute("data-s-section");
                    const a = ue.item[0];
                    const uid = a.dataset.uid;
                    const pN = a.indexOfMe() + 1;
                    api.M[c].update(uid, "position", pN);
                },
            });
        }
        const w = this.children[0].querySelector("a.new");
        if (w !== null) {
            w.addEventListener("click", async () => {
                const c = this.parentElement.getAttribute("data-s-section");
                if (w.classList.contains("skip")) {
                    return api.M[c].create();
                }
                const { value: name } = await Swal({
                    title: "Enter the new role's name",
                    input: "text",
                    showCancelButton: true,
                    inputValidator: (value) => !value && "You need to write something!",
                });
                if (name === undefined) return;
                return api.M[c].create(name);
            });
        }
    }

    get count() {
        let n = 0;
        for (const item of this.children[0].children) {
            if (item.dataset.uid !== undefined) {
                n += 1;
            }
        }
        return n;
    }

    addItem(item) {
        const rlist = this.children[0];
        const oLen = rlist.children.length;
        const nEl = create_element("a", [["data-uid", item.uuid]], [dcTN(item.name)]);
        nEl.addEventListener("click", (e) => {
            const et = e.target;
            this.setActive(Array.from(et.parentElement.children).indexOf(et));
        });
        rlist.insertBefore(
            nEl,
            rlist.querySelector(".div"),
        );
        if (oLen === 2) {
            rlist.parentElement.classList.add("active");
            this.setActive(0);
        }
        const iLen = this.children[0].querySelectorAll("a[data-uid]").length;
        this.setActive(iLen - 1);
    }

    setActive(i) {
        const rlist = this.children[0];
        deActivateChild(rlist);
        rlist.children[i].classList.add("active");
        const c = this.parentElement.getAttribute("data-s-section");
        const r = api.M[c].get(rlist.children[i].dataset.uid);
        const tin = rlist.parentElement.querySelectorAll("[fill]");
        for (const item of tin) {
            item.setAttribute("fill", r.uuid);
            const n = item.getAttribute("name");
            const m = rlist.parentElement.querySelector(`[name="${n}"]`);
            if (m === null) continue;
            m.setAttribute("value", r[n]);
        }
    }

    removeItem(uid) {
        let r = false;
        for (let i = 0; i < this.children[0].children.length; i++) {
            const element = this.children[0].children[i];
            if (element.dataset.uid === uid) {
                if (element.classList.contains("active")) {
                    r = true;
                }
                element.remove();
            }
        }
        if (r) {
            if (this.count === 0) {
                this.classList.remove("active");
                return;
            }
            this.setActive(this.count - 1);
        }
    }

    items() {
        const n = [];
        for (const item of this.children[0].children) {
            n.push(item.dataset.uid);
        }
        return n;
    }
});
