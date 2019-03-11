SCHEDULE = '@every 3s'

filepath = require("filepath")

package.path = filepath.dir(debug.getinfo(1).source)..'/../?.lua;'.. package.path
settings = require "settings"

sendTelegram("Hello from lua!")




