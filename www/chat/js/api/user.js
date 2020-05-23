"use strict";
//
import * as api from "./index.js";

//
/** @type {Map<string,User} */
export const cache = new Map();

//
export class User {
    constructor(o) {
        if (o === null) {
            this.is_null = true;
            return;
        }
        Object.assign(this, o);
        this.is_null = false;
        cache.set(this.uuid, this);
    }
    /** @returns {Promise<api.Role[]>} */
    async getRoles() {
        return Promise.all(this.roles.map((v) => api.M.roles.get(v))).then((l) => {
            return l.sort((a,b) => a.position - b.position);
        });
    }
    /** @returns {string} */
    getName() {
        if (this.nickname.length > 0) {
            return this.nickname;
        }
        return this.name;
    }
    /** @returns {Promise<string>} */
    async getHightestDistinguishedRoleUID() {
        const o = await api.M.users.get(this.uuid);
        const r = await o.getRoles();
        const l = r.filter((v) => v.distinguish);
        const d = l.length > 0 ? l[0].uuid : "";
        return d;
    }
    /** @returns {Promise<string>} */
    async getHightestColoredRoleUID() {
        const o = await api.M.users.get(this.uuid);
        const r = await o.getRoles();
        const l = r.filter((v) => v.color.length > 0);
        const d = l.length > 0 ? l[0].uuid : "";
        return d;
    }
}
