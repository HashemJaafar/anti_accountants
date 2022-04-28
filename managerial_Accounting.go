package main

import "log"

func CALCULATE_CVP_MAP(cvp map[string]map[string]float64, print, check_if_keys_in_the_equations bool) {
	for _, i := range cvp {
		COST_VOLUME_PROFIT(print, check_if_keys_in_the_equations, i)
		_, ok_variable_cost_per_units := i["variable_cost_per_units"]
		if !ok_variable_cost_per_units {
			i["variable_cost_per_units"] = 0
			COST_VOLUME_PROFIT(print, false, i)
		}
		_, ok_fixed_cost := i["fixed_cost"]
		if !ok_fixed_cost {
			i["fixed_cost"] = 0
			COST_VOLUME_PROFIT(print, false, i)
		}
		_, ok_sales_per_units := i["sales_per_units"]
		if !ok_sales_per_units {
			i["sales_per_units"] = 0
			COST_VOLUME_PROFIT(print, false, i)
		}
		_, ok_units := i["units"]
		if !ok_units {
			i["units"] = 0
			COST_VOLUME_PROFIT(print, false, i)
		}
	}
}

func COST_VOLUME_PROFIT(print, check_if_keys_in_the_equations bool, m map[string]float64) {
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
	EQUATIONS_SOLVER(print, check_if_keys_in_the_equations, m, equations)
}

func MIX_COST_VOLUME_PROFIT(print, check_if_keys_in_the_equations bool, m map[string]map[string]float64) {
	var units, sales, variable_cost, fixed_cost float64
	for _, v1 := range m {
		units += v1["units"]
		sales += v1["sales"]
		variable_cost += v1["variable_cost"]
		fixed_cost += v1["fixed_cost"]
	}
	m["total"] = map[string]float64{"units": units, "sales": sales, "variable_cost": variable_cost, "fixed_cost": fixed_cost}
	COST_VOLUME_PROFIT(print, false, m["total"])
}

func COST_VOLUME_PROFIT_SLICE(cvp map[string]map[string]float64, distribution_steps []ONE_STEP_DISTRIBUTION, print, simple bool) {
	CALCULATE_CVP_MAP(cvp, print, true)
	for _, step := range distribution_steps {
		var total_mixed_cost, total_portions_to, total_column_to_distribute float64
		if simple {
			total_mixed_cost = step.AMOUNT
		} else {
			total_mixed_cost = TOTAL_MIXED_COST_IN_COMPLICATED_AND_MULTI_LEVEL_STEP(cvp, step, total_mixed_cost)
		}
		for key_portions_to, portions_to := range step.TO {
			total_portions_to += portions_to
			total_column_to_distribute += cvp[key_portions_to][step.DISTRIBUTION_METHOD]
		}
		for key_portions_to, portions_to := range step.TO {
			var total_overhead_cost_to_sum float64
			switch step.DISTRIBUTION_METHOD {
			case "units_gap":
				total_overhead_cost_to_sum = cvp[key_portions_to]["units_gap"] * cvp[key_portions_to]["variable_cost_per_units"]
				cvp[key_portions_to]["units"] -= cvp[key_portions_to]["units_gap"]
				cvp[key_portions_to]["units_gap"] = 0
			case "1":
				total_overhead_cost_to_sum = total_mixed_cost
			case "portions":
				total_overhead_cost_to_sum = portions_to / total_portions_to * total_mixed_cost
			case "units":
				total_overhead_cost_to_sum = cvp[key_portions_to]["units"] / total_column_to_distribute * total_mixed_cost
			case "variable_cost":
				total_overhead_cost_to_sum = cvp[key_portions_to]["variable_cost"] / total_column_to_distribute * total_mixed_cost
			case "fixed_cost":
				total_overhead_cost_to_sum = cvp[key_portions_to]["fixed_cost"] / total_column_to_distribute * total_mixed_cost
			case "mixed_cost":
				total_overhead_cost_to_sum = cvp[key_portions_to]["mixed_cost"] / total_column_to_distribute * total_mixed_cost
			case "sales":
				total_overhead_cost_to_sum = cvp[key_portions_to]["sales"] / total_column_to_distribute * total_mixed_cost
			case "profit":
				total_overhead_cost_to_sum = cvp[key_portions_to]["profit"] / total_column_to_distribute * total_mixed_cost
			case "contribution_margin":
				total_overhead_cost_to_sum = cvp[key_portions_to]["contribution_margin"] / total_column_to_distribute * total_mixed_cost
			case "percent_from_variable_cost":
				total_overhead_cost_to_sum = cvp[key_portions_to]["variable_cost"] * portions_to
			case "percent_from_fixed_cost":
				total_overhead_cost_to_sum = cvp[key_portions_to]["fixed_cost"] * portions_to
			case "percent_from_mixed_cost":
				total_overhead_cost_to_sum = cvp[key_portions_to]["mixed_cost"] * portions_to
			case "percent_from_sales":
				total_overhead_cost_to_sum = cvp[key_portions_to]["sales"] * portions_to
			case "percent_from_profit":
				total_overhead_cost_to_sum = cvp[key_portions_to]["profit"] * portions_to
			case "percent_from_contribution_margin":
				total_overhead_cost_to_sum = cvp[key_portions_to]["contribution_margin"] * portions_to
			default:
				total_overhead_cost_to_sum = float64(len(step.TO)) * total_mixed_cost
			}
			switch step.SALES_OR_VARIABLE_OR_FIXED {
			case "variable_cost":
				cvp[key_portions_to]["variable_cost_per_units"] = ((cvp[key_portions_to]["variable_cost_per_units"] * cvp[key_portions_to]["units"]) + total_overhead_cost_to_sum) / cvp[key_portions_to]["units"]
			case "fixed_cost":
				cvp[key_portions_to]["fixed_cost"] += total_overhead_cost_to_sum
			default:
				cvp[key_portions_to]["sales_per_units"] = ((cvp[key_portions_to]["sales_per_units"] * cvp[key_portions_to]["units"]) - total_overhead_cost_to_sum) / cvp[key_portions_to]["units"]
			}
			for key_name, map_cvp := range cvp {
				cvp[key_name] = map[string]float64{"units_gap": map_cvp["units_gap"], "units": map_cvp["units"],
					"sales_per_units": map_cvp["sales_per_units"], "variable_cost_per_units": map_cvp["variable_cost_per_units"], "fixed_cost": map_cvp["fixed_cost"]}
			}
			CALCULATE_CVP_MAP(cvp, print, false)
		}
	}
	TOTAL_COST_VOLUME_PROFIT(cvp, print)
}

func LABOR_COST(print, check_if_keys_in_the_equations bool, m map[string]float64) {
	equations := [][]string{
		{"overtime_wage_rate", "bonus_percentage", "*", "hourly_wage_rate"},

		{"overtime_wage", "overtime_hours", "*", "overtime_wage_rate"},
		{"work_time_wage", "work_time_hours", "*", "hourly_wage_rate"},
		{"holiday_wage", "holiday_hours", "*", "hourly_wage_rate"},
		{"vacations_wage", "vacations_hours", "*", "hourly_wage_rate"},
		{"normal_lost_time_wage", "normal_lost_time_hours", "*", "hourly_wage_rate"},
		{"abnormal_lost_time_wage", "abnormal_lost_time_hours", "*", "hourly_wage_rate"},
	}
	EQUATIONS_SOLVER(print, check_if_keys_in_the_equations, m, equations)
}

func PROCESS_COSTING(print, check_if_keys_in_the_equations bool, m map[string]float64) {
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
	EQUATIONS_SOLVER(print, check_if_keys_in_the_equations, m, equations)
}

func TOTAL_COST_VOLUME_PROFIT(cvp map[string]map[string]float64, print bool) {
	var units, sales, variable_cost, fixed_cost float64
	for key_name, map_name := range cvp {
		if key_name != "total" {
			units += map_name["units"]
			sales += map_name["sales"]
			variable_cost += map_name["variable_cost"]
			fixed_cost += map_name["fixed_cost"]
		}
	}
	cvp["total"] = map[string]float64{"units": units, "sales": sales, "variable_cost": variable_cost, "fixed_cost": fixed_cost}
	COST_VOLUME_PROFIT(print, false, cvp["total"])
}

func TOTAL_MIXED_COST_IN_COMPLICATED_AND_MULTI_LEVEL_STEP(cvp map[string]map[string]float64, step ONE_STEP_DISTRIBUTION, total_mixed_cost float64) float64 {
	for key_portions_from, portions := range step.FROM {
		if cvp[key_portions_from]["units"] < portions {
			log.Panic(portions, " for ", key_portions_from, " in ", step.FROM, " is more than ", cvp[key_portions_from]["units"])
		}
		total_mixed_cost += portions * cvp[key_portions_from]["mixed_cost_per_units"]
		cvp[key_portions_from]["fixed_cost"] -= (cvp[key_portions_from]["fixed_cost"] / cvp[key_portions_from]["units"]) * portions
		cvp[key_portions_from]["units"] -= portions
		if cvp[key_portions_from]["units"] == 0 {
			cvp[key_portions_from]["variable_cost_per_units"] = 0
		}
	}
	return total_mixed_cost
}
