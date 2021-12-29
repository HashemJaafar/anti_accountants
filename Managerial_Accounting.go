package anti_accountants

import "log"

type ONE_STEP_DISTRIBUTION struct {
	SALES_OR_VARIABLE_OR_FIXED, DISTRIBUTION_METHOD string
	FROM, TO                                        map[string]float64
}

type MANAGERIAL_ACCOUNTING struct {
	CVP                map[string]map[string]float64
	DISTRIBUTION_STEPS []ONE_STEP_DISTRIBUTION
	PRINT              bool
}

func (s MANAGERIAL_ACCOUNTING) COST_VOLUME_PROFIT_SLICE() {
	s.calculate_cvp_map()
	for _, step := range s.DISTRIBUTION_STEPS {
		var total_mixed_cost, total_portions_to, total_column_to_distribute float64
		for key_portions_from, portions := range step.FROM {
			if s.CVP[key_portions_from]["units"] < portions {
				log.Panic(portions, " for ", key_portions_from, " in ", step.FROM, " is more than ", s.CVP[key_portions_from]["units"])
			}
			total_mixed_cost += portions * s.CVP[key_portions_from]["mixed_cost_per_units"]
			s.CVP[key_portions_from]["fixed_cost"] -= (s.CVP[key_portions_from]["fixed_cost"] / s.CVP[key_portions_from]["units"]) * portions
			s.CVP[key_portions_from]["units"] -= portions
			if s.CVP[key_portions_from]["units"] == 0 {
				s.CVP[key_portions_from]["variable_cost_per_units"] = 0
			}
		}
		for key_portions_to, portions_to := range step.TO {
			total_portions_to += portions_to
			total_column_to_distribute += s.CVP[key_portions_to][step.DISTRIBUTION_METHOD]
		}
		for key_portions_to, portions_to := range step.TO {
			var total_overhead_cost_to_sum float64
			switch step.DISTRIBUTION_METHOD {
			case "units_gap":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["units_gap"] * s.CVP[key_portions_to]["variable_cost_per_units"]
				s.CVP[key_portions_to]["units"] -= s.CVP[key_portions_to]["units_gap"]
				s.CVP[key_portions_to]["units_gap"] = 0
			case "1":
				total_overhead_cost_to_sum = total_mixed_cost
			case "equally":
				total_overhead_cost_to_sum = float64(len(step.TO)) * total_mixed_cost
			case "portions":
				total_overhead_cost_to_sum = portions_to / total_portions_to * total_mixed_cost
			case "units":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["units"] / total_column_to_distribute * total_mixed_cost
			case "variable_cost":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["variable_cost"] / total_column_to_distribute * total_mixed_cost
			case "fixed_cost":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["fixed_cost"] / total_column_to_distribute * total_mixed_cost
			case "mixed_cost":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["mixed_cost"] / total_column_to_distribute * total_mixed_cost
			case "sales":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["sales"] / total_column_to_distribute * total_mixed_cost
			case "profit":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["profit"] / total_column_to_distribute * total_mixed_cost
			case "contribution_margin":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["contribution_margin"] / total_column_to_distribute * total_mixed_cost
			case "percent_from_variable_cost":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["variable_cost"] * portions_to
			case "percent_from_fixed_cost":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["fixed_cost"] * portions_to
			case "percent_from_mixed_cost":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["mixed_cost"] * portions_to
			case "percent_from_sales":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["sales"] * portions_to
			case "percent_from_profit":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["profit"] * portions_to
			case "percent_from_contribution_margin":
				total_overhead_cost_to_sum = s.CVP[key_portions_to]["contribution_margin"] * portions_to
			default:
				log.Panic(step.DISTRIBUTION_METHOD, " is not in [units_gap,1,equally,portions,units,variable_cost,fixed_cost,mixed_cost,sales,profit,contribution_margin,percent_from_variable_cost,percent_from_fixed_cost,percent_from_mixed_cost,percent_from_sales,percent_from_profit,percent_from_contribution_margin]")
			}
			switch step.SALES_OR_VARIABLE_OR_FIXED {
			case "sales":
				s.CVP[key_portions_to]["sales_per_units"] = ((s.CVP[key_portions_to]["sales_per_units"] * s.CVP[key_portions_to]["units"]) - total_overhead_cost_to_sum) / s.CVP[key_portions_to]["units"]
			case "variable_cost":
				s.CVP[key_portions_to]["variable_cost_per_units"] = ((s.CVP[key_portions_to]["variable_cost_per_units"] * s.CVP[key_portions_to]["units"]) + total_overhead_cost_to_sum) / s.CVP[key_portions_to]["units"]
			case "fixed_cost":
				s.CVP[key_portions_to]["fixed_cost"] += total_overhead_cost_to_sum
			default:
				log.Panic(step.SALES_OR_VARIABLE_OR_FIXED, " is not in [sales,variable_cost,fixed_cost]")
			}
			for key_name, map_cvp := range s.CVP {
				s.CVP[key_name] = map[string]float64{"units_gap": map_cvp["units_gap"], "units": map_cvp["units"],
					"sales_per_units": map_cvp["sales_per_units"], "variable_cost_per_units": map_cvp["variable_cost_per_units"], "fixed_cost": map_cvp["fixed_cost"]}
			}
			s.calculate_cvp_map()
		}
	}
	s.total_cost_volume_profit()
}

func (s MANAGERIAL_ACCOUNTING) calculate_cvp_map() {
	for _, i := range s.CVP {
		s.cost_volume_profit(i)
		_, ok_variable_cost_per_units := i["variable_cost_per_units"]
		if !ok_variable_cost_per_units {
			i["variable_cost_per_units"] = 0
			s.cost_volume_profit(i)
		}
		_, ok_fixed_cost := i["fixed_cost"]
		if !ok_fixed_cost {
			i["fixed_cost"] = 0
			s.cost_volume_profit(i)
		}
		_, ok_sales_per_units := i["sales_per_units"]
		if !ok_sales_per_units {
			i["sales_per_units"] = 0
			s.cost_volume_profit(i)
		}
		_, ok_units := i["units"]
		if !ok_units {
			i["units"] = 0
			s.cost_volume_profit(i)
		}
	}
}

func (s MANAGERIAL_ACCOUNTING) total_cost_volume_profit() {
	var units, sales, variable_cost, fixed_cost float64
	for key_name, map_name := range s.CVP {
		if key_name != "total" {
			units += map_name["units"]
			sales += map_name["sales"]
			variable_cost += map_name["variable_cost"]
			fixed_cost += map_name["fixed_cost"]
		}
	}
	s.CVP["total"] = map[string]float64{"units": units, "sales": sales, "variable_cost": variable_cost, "fixed_cost": fixed_cost}
	s.cost_volume_profit(s.CVP["total"])
}

func (s MANAGERIAL_ACCOUNTING) cost_volume_profit(m map[string]float64) {
	equations := [][]string{
		{"variable_cost", "variable_cost_per_units", "*", "units"},
		{"fixed_cost", "fixed_cost_per_units", "*", "units"},
		{"mixed_cost", "mixed_cost_per_units", "*", "units"},
		{"sales", "sales_per_units", "*", "units"},
		{"profit", "profit_per_units", "*", "units"},
		{"contribution_margin", "contribution_margin_per_units", "*", "units"},
		{"mixed_cost", "fixed_cost", "+", "variable_cost"},
		{"sales", "profit", "+", "mixed_cost"},
		{"contribution_margin", "sales", "-", "variable_cost"},
		{"break_even_in_sales", "break_even_in_units", "*", "sales_per_units"},
		{"break_even_in_units", "contribution_margin_per_units", "/", "fixed_cost"},
		{"contribution_margin_per_units", "contribution_margin_ratio", "*", "sales_per_units"},
		{"contribution_margin", "degree_of_operating_leverage", "*", "profit"},
		{"units_gap", "units", "-", "actual_units"},
	}
	check_map_keys_for_equations(equations, m)
	EQUATIONS_SOLVER(s.PRINT, m, equations)
}

func (s MANAGERIAL_ACCOUNTING) PROCESS_COSTING(m map[string]float64) {
	equations := [][]string{
		{"increase_or_decrease", "increase", "-", "decrease"},
		{"increase_or_decrease", "ending_balance", "-", "beginning_balance"},
		{"cost_of_goods_sold", "decrease", "-", "decreases_in_account_caused_by_not_sell"},
		{"equivalent_units", "number_of_partially_completed_units", "*", "percentage_completion"},
		{"equivalent_units_of_production_weighted_average_method", "units_transferred_to_the_next_department_or_to_finished_goods", "+", "equivalent_units_in_ending_work_in_process_inventory"},
		{"cost_of_ending_work_in_process_inventory", "cost_of_beginning_work_in_process_inventory", "+", "cost_added_during_the_period"},
		{"cost_per_equivalent_unit_weighted_average_method", "cost_of_ending_work_in_process_inventory", "/", "equivalent_units_of_production_weighted_average_method"},
		{"equivalent_units_of_production_fifo_method", "equivalent_units_of_production_weighted_average_method", "-", "equivalent_units_in_beginning_work_in_process_inventory"},
		{"percentage_completion_minus_one", "1", "-", "percentage_completion"},
		{"equivalent_units_to_complete_beginning_work_in_process_inventory", "equivalent_units_in_beginning_work_in_process_inventory", "*", "percentage_completion_minus_one"},
		{"cost_per_equivalent_unit_fifo_method", "cost_added_during_the_period", "/", "equivalent_units_of_production_fifo_method"},
	}
	check_map_keys_for_equations(equations, m)
	EQUATIONS_SOLVER(s.PRINT, m, equations)
}

func check_map_keys_for_equations(equations [][]string, m map[string]float64) {
	var elements []string
	for _, equation := range equations {
		elements = append(elements, equation[0], equation[1], equation[3])
	}
	for keyb := range m {
		if !IS_IN(keyb, elements) {
			log.Panic(keyb, " is not in ", elements)
		}
	}
}
