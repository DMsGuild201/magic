### Overview
build an open source mgto gather data scraper. It needs to support multiple data outputs. [Gatherer Extractor](http://www.mtgsalvation.com/forums/magic-fundamentals/magic-software/337224-mtg-gatherer-extractor-v4-0-database-pics) seems to be a standard people are liking. I'd like to change this. The format for Gatherer extractor isn't data friendly. I'm aiming to support it and improve upon it.

### Project status
The project is currently in development. This is NOT a stable API. We're aiming for a stable interface that everyone can follow, and we can continuely patch as MTGO changes.

Right now we're saving the data in JSON format (sets and cards).

The previous way to parse the card's was to match on css. That wont work. Based on the missing data, the rows wont match up. The new strategy is to index all the rows and their values. Then match on the name of the data. It's a rough draft right now, it will be cleaned up as soon as I don't find anymore issues.

TODOO
* Image extraction and saving
-- Support multiple backends (S3, filesystem, etc..)
* Unit test the existing card types
* Pull translations of the cards
* Process card descriptions to support symbols (tap, mana cost etc..)
* Double check we have an interface for ever interaction.


### Magic card edge cases (so far)
Reaper King
http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=159408
{2W}

Westvale Abbey
http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=410049

Demigod of Revenge
http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=370463

Rune-Tail, Kitsune Ascendant
http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=87600

Jace, the Mind Sculptor
http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=413599
