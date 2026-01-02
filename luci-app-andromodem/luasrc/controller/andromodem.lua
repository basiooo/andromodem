module("luci.controller.andromodem", package.seeall)

function index()
    entry({"admin", "modem", "andromodem"}, template("andromodem"), _("Andromodem"), 4).leaf = true
end
