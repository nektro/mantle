"use strict";
//
import * as ui from "./ui.js";
import * as api from "./api/index.js";
import { Channel } from "./ui.channel.js";
import { output } from "./ui.util.js";
import { setDataBinding } from "./util.js";

//
export const M = {
    pong: () => {
        // do nothing, keep connection alive
    },
    user: {
        connect: (d) => {
            ui.M.user.connect(d.user);
        },
        disconnect: (d) => {
            ui.M.user.disconnect(d.user);
        },
        update: (d) => {
            new api.User(d.user);
            if (["add_role","remove_role"].includes(d.key)) {
                document.querySelectorAll(`dialog.popup.user ol [data-role="${d.value}"]`).forEach((v) => {
                    v.classList.toggle("active");
                });
            }
        },
    },
    channel: {
        new: (d) => {
            ui.addChannel(new api.Channel(d.channel));
        },
        update: (d) => {
            const c = new api.Channel(d.channel);
            if (["name"].includes(d.key)) {
                new Channel(c.uuid).p_name = d.value;
            }
            if (output.active_channel_uid === c.uuid) {
                setDataBinding("channel_name", c.name);
                setDataBinding("channel_description", c.description);
            }
        },
    },
    role: {
        new: (d) => {
            ui.M.role.add(new api.Role(d.role));
        },
        update: (d) => {
            new api.Role(d.role);
            if (["color"].includes(d.key)) {
                const x = document.getElementById("link-role-color");
                const y = x.href.split("=");
                x.href = y[0]+"="+(parseInt(y[1],10)+1).toString();
            }
        },
    },
    message: {
        new: async (d) => {
            const u = await api.M.users.get(d.message.author);
            const m = new api.Message(d.message, d.in);
            await output.addMessage(d.in, u, m);
        },
        delete: (d) => {
            for (const item of d.affected) {
                api.M.channels.with(d.channel).messages.remove(item);
            }
        },
    },
    invite: {
        new: (d) => {
            ui.M.invite.add(new api.Invite(d.invite));
        },
        update: (d) => {
            new api.Invite(d.invite);
        },
        delete: (d) => {
            api.M.invites.remove(d.invite);
            ui.M.invite.remove(d.invite);
        },
    },
};
