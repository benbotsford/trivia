#!/usr/bin/env bash
# seed_questions.sh — seeds 10 questions into each question bank.
#
# Usage (from apps/api/):
#   chmod +x seed_questions.sh && ./seed_questions.sh
#
# Reads DEV_AUTH_TOKEN from .env automatically.

set -euo pipefail

# ---- config ------------------------------------------------------------------
if [[ -f .env ]]; then
  set -a; source .env; set +a
fi

TOKEN="${DEV_AUTH_TOKEN:-quibble-dev-token}"
BASE="${API_BASE:-http://localhost:8080}"
AUTH="Authorization: Bearer $TOKEN"

echo "→ Using API: $BASE"
echo "→ Token:     ${TOKEN:0:12}..."
echo ""

# ---- helpers -----------------------------------------------------------------

# json_str <value>  →  JSON-encodes a plain string (handles apostrophes etc.)
json_str() {
  printf '%s' "$1" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))'
}

post_question() {
  local bank_id="$1"
  local body="$2"
  local resp
  resp=$(curl -sf -o /dev/null -w "%{http_code}" \
    -X POST "$BASE/banks/$bank_id/questions" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d "$body")
  if [[ "$resp" != "201" ]]; then
    echo "    ⚠ HTTP $resp posting question" >&2
  fi
}

# add_text <bank_id> <prompt> <answer1> [answer2 ...]
add_text() {
  local bank_id="$1"; local prompt="$2"; shift 2
  local answers_json=""
  for a in "$@"; do
    answers_json+="$(json_str "$a"),"
  done
  answers_json="[${answers_json%,}]"
  post_question "$bank_id" \
    "{\"type\":\"text\",\"prompt\":$(json_str "$prompt"),\"accepted_answers\":$answers_json}"
}

# add_mc <bank_id> <prompt> <correct> <wrong1> <wrong2> <wrong3>
add_mc() {
  local bank_id="$1"; local prompt="$2"; local correct="$3"; shift 3
  local choices_json
  choices_json="{\"text\":$(json_str "$correct"),\"correct\":true}"
  for wrong in "$@"; do
    choices_json+=",{\"text\":$(json_str "$wrong"),\"correct\":false}"
  done
  post_question "$bank_id" \
    "{\"type\":\"multiple_choice\",\"prompt\":$(json_str "$prompt"),\"choices\":[$choices_json]}"
}

# ---- question sets -----------------------------------------------------------

seed_manga_anime() {
  local id="$1"
  add_mc  "$id" "Which Studio Ghibli film follows a girl named Chihiro working in a spirit bathhouse?" \
    "Spirited Away" "My Neighbor Totoro" "Princess Mononoke" "Howl's Moving Castle"
  add_text "$id" "What is the name of the pirate crew led by Monkey D. Luffy in One Piece?" \
    "Straw Hat Pirates" "Straw Hats"
  add_mc  "$id" "In Attack on Titan, what is the name of the special ops squad that protects Eren?" \
    "Survey Corps" "Military Police" "Garrison Regiment" "Titan Hunters"
  add_text "$id" "Who is the author of the manga Naruto?" \
    "Masashi Kishimoto" "Kishimoto"
  add_mc  "$id" "In Dragon Ball Z, what transformation turns a Saiyan's hair golden and multiplies their power?" \
    "Super Saiyan" "Kaio-ken" "Ultra Instinct" "Great Ape"
  add_text "$id" "What is the name of the high school attended by students in My Hero Academia?" \
    "UA High School" "U.A. High" "UA"
  add_mc  "$id" "In Death Note, what must you know to kill someone using the notebook?" \
    "Their face and name" "Their name only" "Their address" "Their date of birth"
  add_text "$id" "What is the name of the underground city in Made in Abyss?" \
    "The Abyss"
  add_mc  "$id" "Which anime features alchemists Edward and Alphonse Elric searching for the Philosopher's Stone?" \
    "Fullmetal Alchemist" "Blue Exorcist" "Soul Eater" "Fairy Tail"
  add_text "$id" "In Demon Slayer, what is the name of the ranking system used by the Demon Slayer Corps?" \
    "Hashira" "The Hashira ranks" "Pillars"
}

seed_arts_literature() {
  local id="$1"
  add_mc  "$id" "Who wrote the novel 'Pride and Prejudice'?" \
    "Jane Austen" "Charlotte Brontë" "George Eliot" "Mary Shelley"
  add_text "$id" "Which Dutch post-impressionist painter cut off part of his own ear?" \
    "Vincent van Gogh" "Van Gogh"
  add_mc  "$id" "In Shakespeare's 'Romeo and Juliet', what are the two feuding family names?" \
    "Montague and Capulet" "Verona and Mantua" "Benvolio and Tybalt" "Capulet and Mercutio"
  add_text "$id" "Who wrote the dystopian novel '1984'?" \
    "George Orwell" "Orwell"
  add_mc  "$id" "Which art movement is Salvador Dalí best associated with?" \
    "Surrealism" "Cubism" "Dadaism" "Expressionism"
  add_text "$id" "Who wrote 'The Great Gatsby'?" \
    "F. Scott Fitzgerald" "Fitzgerald" "F Scott Fitzgerald"
  add_mc  "$id" "The Sistine Chapel ceiling was painted by which Renaissance artist?" \
    "Michelangelo" "Leonardo da Vinci" "Raphael" "Donatello"
  add_text "$id" "Which American poet wrote 'The Raven' in 1845?" \
    "Edgar Allan Poe" "Poe"
  add_mc  "$id" "Which novel begins with the line 'Call me Ishmael'?" \
    "Moby-Dick" "The Old Man and the Sea" "Billy Budd" "Lord Jim"
  add_text "$id" "In which city is the Louvre Museum located?" \
    "Paris"
}

seed_science_nature() {
  local id="$1"
  add_mc  "$id" "What is the chemical symbol for gold?" \
    "Au" "Go" "Gd" "Ag"
  add_text "$id" "What planet is known as the Red Planet?" \
    "Mars"
  add_mc  "$id" "What is the powerhouse of the cell?" \
    "Mitochondria" "Nucleus" "Ribosome" "Golgi apparatus"
  add_text "$id" "What is the most abundant gas in Earth's atmosphere?" \
    "Nitrogen"
  add_mc  "$id" "At what temperature (°C) does water boil at sea level?" \
    "100" "90" "110" "80"
  add_text "$id" "How many bones are in the adult human body?" \
    "206"
  add_mc  "$id" "Which planet in our solar system has the most moons?" \
    "Saturn" "Jupiter" "Uranus" "Neptune"
  add_text "$id" "What force keeps planets in orbit around the Sun?" \
    "Gravity" "Gravitational force"
  add_mc  "$id" "What is the atomic number of carbon?" \
    "6" "4" "8" "12"
  add_text "$id" "What is the speed of light in a vacuum, approximately in km/s?" \
    "300000" "299792" "300,000"
}

seed_general_knowledge() {
  local id="$1"
  add_mc  "$id" "How many time zones does Russia span?" \
    "11" "9" "13" "7"
  add_text "$id" "What is the currency of Japan?" \
    "Yen" "Japanese Yen"
  add_mc  "$id" "Which country is home to the Great Barrier Reef?" \
    "Australia" "Indonesia" "Philippines" "Fiji"
  add_text "$id" "How many sides does a pentagon have?" \
    "5" "Five"
  add_mc  "$id" "What is the most widely spoken language in the world by total speakers?" \
    "English" "Mandarin Chinese" "Spanish" "Hindi"
  add_text "$id" "What is the smallest country in the world by area?" \
    "Vatican City" "Vatican"
  add_mc  "$id" "In what year did the first moon landing take place?" \
    "1969" "1965" "1971" "1968"
  add_text "$id" "What is the chemical symbol for iron?" \
    "Fe"
  add_mc  "$id" "Which ocean is the deepest in the world?" \
    "Pacific Ocean" "Atlantic Ocean" "Indian Ocean" "Arctic Ocean"
  add_text "$id" "How many players are on a standard chess team in team competitions?" \
    "4"
}

# ---- dispatch by bank name ---------------------------------------------------

seed_for_bank() {
  local bank_id="$1"
  local bank_name="$2"
  local name_lc
  name_lc=$(echo "$bank_name" | tr '[:upper:]' '[:lower:]')

  printf "  Seeding %-28s" "'$bank_name'..."

  if echo "$name_lc" | grep -qE 'manga|anime'; then
    seed_manga_anime "$bank_id"
  elif echo "$name_lc" | grep -qE 'art|liter'; then
    seed_arts_literature "$bank_id"
  elif echo "$name_lc" | grep -qE 'sci|nature'; then
    seed_science_nature "$bank_id"
  elif echo "$name_lc" | grep -qE 'general|knowledge'; then
    seed_general_knowledge "$bank_id"
  else
    echo ""
    echo "    ⚠ No matching category for '$bank_name' — skipping"
    return
  fi

  echo " ✓"
}

# ---- fetch banks and run -----------------------------------------------------

banks_json=$(curl -sf -H "$AUTH" "$BASE/banks")
bank_count=$(echo "$banks_json" | python3 -c 'import json,sys; print(len(json.load(sys.stdin)))')

if [[ "$bank_count" -eq 0 ]]; then
  echo "No banks found. Create a bank first via the app, then re-run."
  exit 1
fi

echo "Found $bank_count bank(s)."
echo ""

while IFS=$'\t' read -r bank_id bank_name; do
  seed_for_bank "$bank_id" "$bank_name"
done < <(echo "$banks_json" | python3 -c '
import json, sys
for b in json.load(sys.stdin):
    print(b["id"] + "\t" + b["name"])
')

echo ""
echo "Done!"
