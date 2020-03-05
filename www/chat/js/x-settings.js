"use strict";
//

class SettingsDialog extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        if (this.dataset.active === undefined) {
            this.setActivePane(0);
        }
        for (const item of this.children[0].children[0].children) {
            item.addEventListener("click", (ev) => {
                const t = ev.target;
                const i = Array.from(t.parentElement.children).indexOf(t);
                this.setActivePane(i);
            });
        }
    }
    /**
     * @param {Number} n
     */
    setActivePane(n) {
        if (this.dataset.active !== undefined) {
            const o = parseInt(this.dataset.active, 10);
            this.children[0].children[0].children[o].classList.remove("active");
            this.children[0].children[1].children[o].classList.remove("active");
        }
        this.dataset.active = n.toString();
        this.children[0].children[0].children[n].classList.add("active");
        this.children[0].children[1].children[n].classList.add("active");
    }
}

customElements.define("x-settings", SettingsDialog);
