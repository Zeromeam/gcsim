# Stellar Swirl account calculator

This directory contains reproducible Stellar Swirl simulations built on top of
gcsim's frame-based combat engine. Character builds and rotations remain config
data, while unusual character and reaction behavior is implemented in dedicated
Go modules. A constellation, weapon, refinement, artifact set, stat line, target,
or rotation can therefore be changed without hard-coding a showcase result.

## Audited results

Results below are from 5,000 iterations against one Level 100 enemy with 10%
base RES. They were generated on 2026-09-03 after the formula and regression-test
audit in this branch.

| Configuration | Mean damage | Duration | Mean DPS | DPS standard deviation |
| --- | ---: | ---: | ---: | ---: |
| Current account, modeled rotation | 5,341,432 | 64.983 s | **82,197** | 2,544 |
| Public Prydwen C1 setup, KQMS approximation | 7,613,207 | 84.400 s | **90,204** | 4,800 |
| Same public setup with Mizuki C2 | 8,371,816 | 84.400 s | **99,192** | 5,492 |

The current-account total is an independent model result, not a calibration to
the supplied live observation of 4,074,642 damage / 62,694 DPS. The observed
number remains a useful comparison, but it is not assumed to be correct: the
exact live action sequence, target behavior, missed hits, and buff indicators
were not supplied. Only the target level, base RES, nominal duration, account
builds, and a KQM/Prydwen-style rotation are represented here.

As one narrow sanity check, the modeled current account averages about 222,589
damage for each Coda finisher. That is close to the supplied 225,576 strongest
Odette hit even though the simulator was not tuned to it. Character damage
breakdowns are not directly comparable because gcsim attributes a combined
reaction Stellar Swirl hit to its Anemo trigger while every contributor's stats
still participate in the damage formula.

The often-repeated 146k/156k sheet claims are not treated as ground truth. Their
complete artifact rolls, refinements, proc assumptions, and hit-by-hit timeline
were not available. The transparent public source gives gear, a 21.1-second
rotation, four-rotation averaging, and a relative 89.2% result for the Sucrose
variant, but no absolute DPS. The two public configs here reconstruct those
published assumptions explicitly; their failure to reproduce 146k/156k is a
finding, not something hidden by a correction factor.

## Configurations

- `current_account.gcsim` contains the exact account snapshot used in this
  project: Mizuki C2/Sacrificial Fragments R5/4VV, Odette C0/Primordial Jade
  Cutter R1/4 Furnace, Cryo Aether C2/Exaiphanes Blade R3/2pc+2pc ATK, and
  Sucrose C6/Wandering Evenstar R4/4VV. No account UID is stored.
- `prydwen_c1_kqms.gcsim` is the optimized KQMS reconstruction of Prydwen's C1
  Sucrose variant.
- `prydwen_c2_kqms.gcsim` changes that reconstruction to Mizuki C2 and rebalances
  its liquid substat rolls.
- The corresponding `_base.gcsim` files preserve the pre-optimization KQMS
  inputs. Refinements not displayed by Prydwen are stated as assumptions in the
  configs rather than silently invented.

## Mechanics implemented for this calculation

- Reaction Stellar Swirl uses separate contributor calculations, ranks them,
  and combines them at 60% / 30% / 5% / 5%.
- The initial Anemo hit uses a 0.75 reaction multiplier. Level 1 and Level 2
  Vortex detonations use 2 and 3. Five additional triggers detonate the Vortex
  early.
- Direct and reaction Stellar damage ignore DEF and regular Elemental/Common
  DMG Bonus, but use the applicable scaling stat, EM bonus, Stellar reaction
  bonuses, CRIT, enemy RES, Base DMG Bonus, additive base damage, and Elevation
  in their distinct formula positions.
- Mizuki's Version 7.0 Witch's Revelation, team EM share, Stellar Skill bonus,
  C1 ordinary/SSW coefficients and extra hits, and C2 buffs/RES reduction are
  implemented.
- Odette C0 includes both Skills, the Dance Double, Marvelous Splendor transfer,
  Burst, Snow Swan's Dream, Stellar Base DMG Bonus, and A4 Elevation.
- Cryo Aether C0-C2 includes Foreign Permafrost's seven resonance bonuses,
  Icepoint, Freezing Ice, Frostglow, Skill crystals, Burst conversion, C1 energy,
  and C2 EM.
- Heart of the Furnace, Exaiphanes Blade R1-R5, Stellar-aware Viridescent
  Venerer, and Sucrose's Stellar Swirl interaction are supported.

## Correctness controls

Regression tests cover the core formula, Vortex stack cap, contributor weighting
and additive-damage placement, Mizuki's 550% C1 Stellar coefficient,
Exaiphanes R3's Traveler-only 42% CRIT DMG and 5 Energy, Odette's full ten-hit
summon, Burst-only reaction buff, and the persistence of Coda-converted Dance
Double attacks.

The simulator-only test suite passes:

```sh
GOCACHE=/tmp/gcsim-go-cache go test ./internal/... ./pkg/...
```

The repository-wide `go test ./...` additionally exercises backend integration
packages. On the local Codex machine it reaches the simulator tests but fails in
three unrelated backend tests because Docker/Mongo is unavailable and a
persistent test user already exists.

## Known limits

- Odette's TCL page currently publishes talent values but no frame findings.
  Her Dance Double hitmarks and Coda DoT timing are therefore provisional live
  measurements and remain the largest timing uncertainty.
- The validated scope is a stationary single-target rotation. The inherited
  Stellar Vortex state is global and has not yet been validated for separated
  multi-target Vortices.
- Only the constellations required by these three configs are complete. Odette
  C1-C6, Cryo Traveler C3-C6, Odette Normal/Charged/Plunging attacks, Cryo Aether
  Plunging attacks, and Stellar-Conduct behavior for these new modules are not
  claimed as implemented.
- The exact live rotation can be added as another config when its full button
  sequence is available. It should be compared with this model, not used to
  overwrite verified formulas.

## Run

```sh
go run ./cmd/gcsim -c configs/stellar_swirl/current_account.gcsim -out current.json -nb
go run ./cmd/gcsim -c configs/stellar_swirl/prydwen_c1_kqms.gcsim -out public-c1.json -nb
go run ./cmd/gcsim -c configs/stellar_swirl/prydwen_c2_kqms.gcsim -out public-c2.json -nb
```

## Sources

- KQM Stellar Glimmer mechanics and formulas:
  <https://feelcrafting.com/misc/stellar/>
- KQM Mizuki TCL: <https://library.keqingmains.com/characters/anemo/yumemizuki-mizuki>
- KQM Odette TCL: <https://library.keqingmains.com/characters/cryo/odette>
- KQM Cryo Traveler TCL: <https://library.keqingmains.com/characters/cryo/traveler-cryo>
- Exaiphanes Blade refinement data: <https://paimon.moe/weapons/exaiphanes_blade>
- Prydwen Mizuki builds, public gear, rotations, and relative comparisons:
  <https://www.prydwen.gg/genshin-impact/characters/yumemizuki-mizuki>
- gcsim upstream: <https://github.com/genshinsim/gcsim>
