"use strict";
//
import { cache as cacheU, User } from "./user.js";
import { cache as cacheC, Channel } from "./channel.js";
import { cache as cacheR, Role } from "./role.js";
import { cache as cacheG, Message } from "./message.js";

//
export {
    User,
    Channel,
    Role,
    Message,
};

//
const caches = [
    cacheU,
    cacheC,
    cacheR,
    cacheG,
];

//
function fetchE(endpoint, method="get", data={}) {
    const body = new FormData();
    for (const k in data) {
        if (!Object.prototype.hasOwnProperty.call(data, k)) continue;
        body.set(k, data[k]);
    }
    const opts = method === "get" ? {} : {method, body};
    return fetch(`./../api${endpoint}`, opts).then((x) => {
        if (x.headers.getSafe("content-type").includes("application/json")) return x.json();
        return x.text();
    }).then((x) => {
        if (typeof x === "string") return;
        if (!x.success) {
            return Promise.reject(new Error(x.message));
        }
        return x.message;
    });
}
function fetchI(endpoint, cl) {
    return fetchE(endpoint).then((x) => {
        if (cl === undefined) return x;
        return new cl(x);
    });
}
function fetchL(endpoint, cl) {
    return fetchE(endpoint).then((x) => {
        return x.map((y) => {
            return new cl(y);
        });
    });
}
function fetchIC(endpoint, cl, cch, key) {
    if (caches[cch].has(key)) return caches[cch].get(key);
    return fetchI(endpoint, cl);
}

//
export const M = {
    meta: {
        about: () => {
            return fetchI("/about");
        },
    },
    users: {
        /** @returns {Promise<User>} */
        get: (uid) => {
            return fetchIC(`/users/${uid}`, User, 0, uid);
        },
        me: () => {
            return fetchE("/users/@me").then((x) => {
                return {
                    user: new User(x.me),
                    perms: x.me,
                };
            });
        },
        /** @returns {Promise<User[]>} */
        online: () => {
            return fetchL("/users/online", User);
        },
        update: (uid,k,v) => {
            return fetchE(`/users/${uid}/update`, "put", { p_name: k, p_value: v, });
        },
    },
    channels: {
        /** @returns {Promise<Channel[]>} */
        me: () => {
            return fetchL("/channels/@me", Channel);
        },
        /** @returns {Promise<Channel>} */
        get: (uid) => {
            return fetchIC(`/channels/${uid}`, Channel, 1, uid);
        },
        create: (n) => {
            return fetchE("/channels/create", "post", { name: n });
        },
        update: (uid,k,v) => {
            return fetchE(`/channels/${uid}/update`, "put", { p_name: k, p_value: v, });
        },
    },
    roles: {
        /** @returns {Promise<Role[]>} */
        me: () => {
            return fetchL("/roles", Role);
        },
        /** @returns {Promise<Role>} */
        get: (uid) => {
            return fetchIC(`/roles/${uid}`, Role, 2, uid);
        },
        create: (n) => {
            return fetchE("/roles/create", "post", { name: n });
        },
        update: (uid,k,v) => {
            return fetchE(`/roles/${uid}/update`, "put", { p_name: k, p_value: v, });
        },
    },
};
