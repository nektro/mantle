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
function fetchE(endpoint) {
    return fetch(`./../api${endpoint}`).then((x) => x.json()).then((x) => {
        if (!x.success) {
            return Promise.reject(new Error(x.message));
        }
        return x.message;
    });
}
function fetchI(endpoint, cl) {
    return fetchE(endpoint).then((x) => {
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
    users: {
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
        online: () => {
            return fetchL("/users/online", User);
        }
    },
    channels: {
        me: () => {
            return fetchL("/channels/@me", Channel);
        },
        get: (uid) => {
            return fetchIC(`/channels/${uid}`, Channel, 1, uid);
        },
    },
    roles: {
        getAll: () => {
            return fetchL("/roles", Role);
        },
        get: (uid) => {
            return fetchIC(`/roles/${uid}`, Role, 2, uid);
        },
    },
};
