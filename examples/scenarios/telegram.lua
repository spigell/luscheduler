SCHEDULE = '00 00 18 * * *'

package.path = filepath.dir(debug.getinfo(1).source)..'/../?.lua;'.. package.path
settings = require "settings"

sendTelegram("Hello from lua!")




