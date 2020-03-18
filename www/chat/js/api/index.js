"use strict";
//
import { User } from "./user.js";
import { Channel } from "./channel.js";
import { Role } from "./role.js";
import { Message } from "./message.js";

//
export {
    User,
    Channel,
    Role,
    Message,
};

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

//
export const M = {
    users: {
        get: (uid) => {
            return fetchI(`/users/${uid}`, User);
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
    roles: {
        get: () => {
            return fetchL("/roles", Role);
        },
    },
    channels: {
        me: () => {
            return fetchL("/channels/@me", Channel);
        },
    },
};
