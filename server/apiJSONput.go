package main

import (
	"mule/jsend"
	"mule/overpower/models"
	"net/http"
)

// /overpower/json/RESOURCE
func (h *Handler) ApiJSONput(w http.ResponseWriter, r *http.Request) {
	if !h.LoggedIn {
		JSONUserError(w, "You must be logged in to use JSON PUT")
		return
	}
	resource, ok := h.StrAt(3)
	if !ok {
		JSONUserError(w, "NO RESOURCE SPECIFIED")
		return
	}
	if len(h.Path) > 4 {
		JSONUserError(w, "OVERSPECIFIED URL: TOO MANY ARUGMENTS")
		return
	}
	switch resource {
	case "launchorders":
		h.apiJSONputLaunchOrders(w, r)
	case "factions":
		h.apiJSONputFactions(w, r)
	case "mapviews":
		h.apiJSONputMapViews(w, r)
	default:
		JSONUserError(w, "UNKNOWN RESOURCE REQUESTED", KV{"resource", resource})
	}
}

func (h *Handler) apiJSONputLaunchOrders(w http.ResponseWriter, r *http.Request) {
	item := &models.LaunchOrder{}
	err := jsend.Read(r, &item)
	if err != nil {
		JSONUserError(w, "Cannot read json into launchorder data")
		return
	}
	_, ok, err := h.Validate(item.GID, item.FID)
	if my, bad := Check(err, "API PUT failure on resource validation", "type", "launchorder", "GID", item.GID, "FID", item.FID); bad {
		JSONServerError(w, my)
		return
	}
	if !ok {
		JSONUserError(w, "You are not authorized for that faction")
		return
	}
	errS, errU := InternalSetLaunchOrder(item.GID, item.FID, item.Size, item.Source, item.Target)
	if my, bad := Check(errS, "API JSON PUT LAUNCHORDER failure on command execution", "item", item); bad {
		JSONServerError(w, my)
		return
	}
	if errU != nil {
		JSONUserError(w, errU.Error())
		return
	}
	JSONSuccess(w, nil)
}

func (h *Handler) apiJSONputFactions(w http.ResponseWriter, r *http.Request) {
	item := &models.Faction{}
	err := jsend.Read(r, &item)
	if err != nil {
		JSONUserError(w, "Cannot read json into faction data")
		return
	}
	_, ok, err := h.Validate(item.GID, item.FID)
	if my, bad := Check(err, "API PUT failure on resource validation", "type", "faction", "GID", item.GID, "FID", item.FID); bad {
		JSONServerError(w, my)
		return
	}
	if !ok {
		JSONUserError(w, "You are not authorized for that faction")
		return
	}
	errS, errU := InternalSetDoneBuffer(item.GID, item.FID, item.DoneBuffer)
	if my, bad := Check(errS, "API JSON PUT FACTION failure on command execution", "item", item); bad {
		JSONServerError(w, my)
		return
	}
	if errU != nil {
		JSONUserError(w, errU.Error())
		return
	}
	JSONSuccess(w, nil)
}

func (h *Handler) apiJSONputMapViews(w http.ResponseWriter, r *http.Request) {
	item := &models.MapView{}
	err := jsend.Read(r, &item)
	if err != nil {
		JSONUserError(w, "Cannot read json into mapview data")
		return
	}
	_, ok, err := h.Validate(item.GID, item.FID)
	if my, bad := Check(err, "API PUT failure on resource validation", "type", "mapview", "GID", item.GID, "FID", item.FID); bad {
		JSONServerError(w, my)
		return
	}
	if !ok {
		JSONUserError(w, "You are not authorized for that faction")
		return
	}
	errS, errU := InternalSetMapCenter(item.GID, item.FID, item.Center)
	if my, bad := Check(errS, "API JSON PUT MAPVIEW failure on command execution", "item", item); bad {
		JSONServerError(w, my)
		return
	}
	if errU != nil {
		JSONUserError(w, errU.Error())
		return
	}
	JSONSuccess(w, nil)
}
