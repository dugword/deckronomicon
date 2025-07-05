import json
import matplotlib.pyplot as plt

def load_summary(filename):
    with open("two/" + filename, 'r') as f:
        return json.load(f)

flush_out_summary = load_summary("flush_out_summary.json")
thrill_of_possibility_summary = load_summary("thrill_of_possibility_summary.json")
pearled_unicorn_summary = load_summary("pearled_unicorn_summary.json")

if flush_out_summary["NumberOfTurns"] != thrill_of_possibility_summary["NumberOfTurns"] or \
   flush_out_summary["NumberOfTurns"] != pearled_unicorn_summary["NumberOfTurns"]:
    raise ValueError("Summaries have different number of turns")

number_of_turns = flush_out_summary["NumberOfTurns"]

turns = list(range(1, number_of_turns + 1))

# Flush Out Summary

plt.fill_between(turns, flush_out_summary["AvgCumulativeCardDrawnPerTurn"], label="Cards Drawn w/ Flush Out", alpha=1.0)
# plt.fill_between(turns, flush_out_summary["AvgCumulativeFlushOutDrawnPerTurn"], label="Flush Out Seen", alpha=1.0)

# Thrill of Possibility Summary

plt.fill_between(turns, thrill_of_possibility_summary["AvgCumulativeCardDrawnPerTurn"], label="Cards Drawn w/ Thrill of Possibility", alpha=1.0)
#plt.fill_between(turns, thrill_of_possibility_summary["AvgCumulativeThrillOfPossibilityDrawnPerTurn"], label="Thrill of Possibility Seen", alpha=1.0)

# Pearled Unicorn Summary

plt.fill_between(turns, pearled_unicorn_summary["AvgCumulativeCardDrawnPerTurn"], label="Cards Drawn Base", alpha=1.0)
#plt.fill_between(turns, pearled_unicorn_summary["AvgCumulativePearledUnicornDrawnPerTurn"], label="Pearled Unicorn Seen", alpha=1.0)

plt.xlabel("Turn")
plt.ylabel("Counts")
plt.xticks(turns)
plt.xlim(1, number_of_turns)

plt.title("Results")
plt.legend()
# plt.grid(True)
plt.show()
