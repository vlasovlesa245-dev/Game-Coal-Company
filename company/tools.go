package company

type Tool struct {
	Title string
	Price int
	IsBought bool
}

type StatisticsTools struct {
	Sold []Tool
	Unsold []Tool
}

func GetTools() map[string]Tool {
	tools := make(map[string]Tool)
	tools["Pickaxes"] = Tool{
		Title: "Pickaxes",
		Price: 3000,
		IsBought: false,
	}
	tools["Ventilation"] = Tool{
		Title: "Ventilation",
		Price: 15000,
		IsBought: false,
	}
	tools["Minecarts"] = Tool{
		Title: "Minecarts",
		Price: 50000,
		IsBought: false,
	}
	return tools
}

func (c *Company) BuyTool(title string) {
	c.mtx.Lock()
	bought := c.Tools[title]
	bought.IsBought = true
	c.Tools[title] = bought
	c.Coal.Ballance -= c.Tools[title].Price
	c.mtx.Unlock()
}

func (c *Company) Get_Statistics_Of_Tools() StatisticsTools {
	Sold := make([]Tool, 0, 3)
	Unsold := make([]Tool, 0, 3)

	for t, v := range c.Tools {
		if ok := c.Tools[t].IsBought; ok {
			Sold = append(Sold, v)
		} else {
			Unsold = append(Unsold, v)
		}
	}
	tools := StatisticsTools{
		Sold: Sold,
		Unsold: Unsold,
	}

	return tools
}