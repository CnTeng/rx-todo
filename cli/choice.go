package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/CnTeng/rx-todo/client"
)

func (c *cli) SelectOne(rs client.ResourceSlice) (int64, error) {
	switch rs := rs.(type) {
	case client.ProjectSlice:
		c.PrintProjects(rs, "Projects", client.None)
	case client.TaskSlice:
		c.PrintTasks(rs, "Tasks", client.None)
	default:
		fmt.Println("Unknown resource type")
	}

	fmt.Print("Enter the index, only one: ")

	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	return rs.GetIDByIndex(index)
}

func parseIndexRange(s string) ([]int, error) {
	idxs := make([]int, 0)

	idxRange := strings.Split(s, "-")
	if len(idxRange) != 2 {
		return nil, fmt.Errorf("invalid range '%s'", s)
	}

	start, err := strconv.Atoi(idxRange[0])
	if err != nil {
		return nil, fmt.Errorf("invalid range '%s': %v", s, err)
	}

	end, err := strconv.Atoi(idxRange[1])
	if err != nil {
		return nil, fmt.Errorf("invalid range '%s': %v", s, err)
	}

	if start > end {
		return nil, fmt.Errorf("invalid range '%s': start > end", s)
	}

	for i := start; i <= end; i++ {
		idxs = append(idxs, i)
	}

	return idxs, nil
}

func (c *cli) SelectMultiple(s client.ResourceSlice) ([]int64, error) {
	fmt.Println("Enter the index, ranges (e.g., 1-5), or type 'all': ")

	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	input = strings.TrimSpace(input)
	if input == "all" {
		ids, _ := s.GetIDsByIndexs(nil)
		return ids, nil
	}

	choices := strings.Split(input, ",")
	ids := make([]int64, 0)
	for _, choice := range choices {
		choice = strings.TrimSpace(choice)

		if strings.Contains(choice, "-") {
			idxs, err := parseIndexRange(choice)
			if err != nil {
				return nil, err
			}

			i, err := s.GetIDsByIndexs(idxs)
			if err != nil {
				return nil, err
			}
			ids = append(ids, i...)

			continue
		}

		i, err := strconv.Atoi(choice)
		if err != nil {
			return nil, fmt.Errorf("invalid ID '%s': %v", choice, err)
		}

		id, err := s.GetIDByIndex(i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
