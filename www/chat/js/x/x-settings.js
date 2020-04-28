"use strict";
//
import { deActivateChild } from "./../util.js";

//
class SettingsDialog extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        if (this.dataset.active === undefined) {
            this.setActivePane(0);
        }
        for (const item of this.kids()) {
            item.addEventListener("click", (ev) => {
                const t = ev.target;
                const i = Array.from(this.kids()).indexOf(t);
                this.setActivePane(i);
            });
        }
        this.addEventListener("click", (e) => {
            if (e.target.localName === "x-settings") {
                e.target.removeAttribute("open");
            }
        });
        const mo = new MutationObserver((muts) => {
            for (const item of muts) {
                if (item.type === "attributes") {
                    if (item.attributeName === "hidden") {
                        if (item.target.classList.contains("active")) {
                            deActivateChild(this.nav());
                            deActivateChild(this.pane());
                        }
                    }
                }
            }
        });
        for (const item of this.kids()) {
            mo.observe(item, {
                attributes: true,
            });
        }
    }
    nav() {
        return this.children[0].children[0];
    }
    pane() {
        return this.children[0].children[1];
    }
    kids() {
        return this.nav().querySelectorAll("a:not(.div)");
    }
    /**
     * @param {Number} n
     */
    setActivePane(n) {
        deActivateChild(this.nav());
        deActivateChild(this.pane());
        this.dataset.active = n.toString();
        this.kids()[n].classList.add("active");
        this.pane().children[n].classList.add("active");
    }
    _open() {
        this.setAttribute("open","");
        for (let i = 0; i < this.nav().children.length; i++) {
            const item = this.nav().children[i];
            if (!item.hasAttribute("hidden")) {
                this.setActivePane(i);
                break;
            }
        }
    }
}

customElements.define("x-settings", SettingsDialog);
