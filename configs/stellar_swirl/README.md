# Stellar Swirl account calculator

This directory contains reproducible Stellar Swirl simulations built on top of
gcsim's frame-based combat engine. Character builds and rotations remain config
data, while unusual character and reaction behavior is implemented in dedicated
Go modules. A constellation, weapon, refinement, artifact set, stat line, target,
or rotation can therefore be changed without hard-coding a showcase result.

## Audited results

Results below are against one Level 100 enemy with 10% base RES. They were
generated on 2026-09-04 after the formula, event-timing, and regression-test
audit in this branch. The account and Prydwen rows use 5,000 iterations; the
longer recorded-showcase audit uses 20,000.

| Configuration | Mean damage | Duration | Mean DPS | Maximum DPS | DPS standard deviation |
| --- | ---: | ---: | ---: | ---: | ---: |
| Current account, modeled rotation | 5,997,150 | 64.983 s | **92,288** | 101,345 | 2,770 |
| Public Prydwen C1 setup, KQMS approximation | 8,028,726 | 84.400 s | **95,127** | 111,237 | 4,937 |
| Same public setup with Mizuki C2 | 8,801,575 | 84.400 s | **104,284** | 121,930 | 5,511 |
| Recorded 168k showcase reconstruction | 19,052,230 | 120.017 s | **158,747** | 167,558 | 3,402 |

“Maximum DPS” is the highest observed trial in each finite batch, not the
expected result. The recorded-showcase comparison below uses the reproducible
sample log for seed `9594297515983362552` from its 20,000-trial run.

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
were not available. The Prydwen configs reconstruct the published assumptions
explicitly; their failure to reproduce 146k/156k is a finding, not something
hidden by a correction factor.

### Recorded 168k showcase comparison

The separate `public_168k_showcase.gcsim` config transcribes the builds and
120.03-second action sequence from Key_Influence_8972's recorded dummy test.
Unlike the Prydwen sheet approximation, the recording exposes enough individual
hits to audit the calculation instead of merely comparing final DPS.

| Metric | Recording | Model's maximum seed | Difference |
| --- | ---: | ---: | ---: |
| Total damage | 20,231,019 | 20,109,811 | -0.60% |
| DPS | 168,551 | 167,558 | -0.59% |
| Strongest hit | 429,816 | 429,995 | +0.04% |
| Mizuki damage | 9,527,569 | 9,520,249 | -0.08% |
| Odette damage | 5,200,067 | 5,461,697 | +5.03% |
| Cryo Traveler damage | 4,358,043 | 4,138,511 | -5.04% |
| Sucrose damage | 1,145,340 | 989,354 | -13.62% |

Visible like-for-like hits independently agree: the opening Traveler Skill is
4,092 in the model versus 4,090 in the recording, Sucrose's first Stellar Swirl
is 6,156 versus 6,153, Traveler's Javelin is 67,378 versus 67,341, and
Traveler's two Charged Attack hits are 124,628/140,043 versus
124,559/139,966. This rules out a broad damage multiplier error.

The reconstructed maximum sample ends three final Javelin hits short when the
120.03-second timer expires. Those hits would add about 202,135 Traveler damage,
putting that character within 0.40% of the recording, but they are not added to
the reported simulator result. The remaining character-level differences mostly
offset one another and can come from sub-frame rotation timing, rounded build-card
stats, and differing ownership of combined Stellar Swirl damage. The total and
like-for-like hit comparisons are therefore more reliable than treating either
character breakdown as an absolute truth.

## Configurations

- `current_account.gcsim` contains the exact account snapshot used in this
  project: Mizuki C2/Sacrificial Fragments R5/4VV, Odette C0/Primordial Jade
  Cutter R1/4 Furnace, Cryo Aether C2/Exaiphanes Blade R3/2pc+2pc ATK, and
  Sucrose C6/Wandering Evenstar R4/4VV. No account UID is stored.
- `prydwen_c1_kqms.gcsim` is the optimized KQMS reconstruction of Prydwen's C1
  Sucrose variant.
- `prydwen_c2_kqms.gcsim` changes that reconstruction to Mizuki C2 and rebalances
  its liquid substat rolls.
- `public_168k_showcase.gcsim` contains the recorded C0 Mizuki showcase's exact
  visible builds, uses Cryo Lumine, and follows its measured action timeline.
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
- Cryo Traveler C0-C2 includes Foreign Permafrost's seven resonance bonuses,
  Icepoint, Freezing Ice, Frostglow, Skill crystals, Burst conversion, C1 energy,
  and C2 EM. Aether and Lumine use their own Normal/Charged frame data and
  second Charged Attack multipliers while sharing the Cryo kit mechanics.
- Heart of the Furnace, Exaiphanes Blade R1-R5, Stellar-aware Viridescent
  Venerer, and Sucrose's Stellar Swirl interaction are supported.

## Correctness controls

Regression tests cover the core formula, Vortex stack cap, contributor weighting
and additive-damage placement, Mizuki's 550% C1 Stellar coefficient,
Exaiphanes R3's Traveler-only 42% CRIT DMG and 5 Energy, Lumine's female
Charged Attack multiplier and N1 timing, Heart of the Furnace's party bonus
regardless of config order, Odette's Coda timing, all ten Dance Double Cryo
applications, Burst-only reaction buff, and the persistence of Coda-converted
Dance Double attacks.

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
  Her Dance Double hitmarks, per-hit Cryo application, and Coda timing are
  measured from the recorded live footage and remain the largest timing
  uncertainty.
- The validated scope is a stationary single-target rotation. The inherited
  Stellar Vortex state is global and has not yet been validated for separated
  multi-target Vortices.
- Only the constellations required by these configs are complete. Odette
  C1-C6, Cryo Traveler C3-C6, Odette Normal/Charged/Plunging attacks, Cryo Aether
  Plunging attacks, and Stellar-Conduct behavior for these new modules are not
  claimed as implemented.
- Video timestamps cannot reveal every internal game frame, rounded stat-card
  decimal, random seed, or ownership rule used by an external damage tracker.
  The showcase config is therefore a measured reconstruction, not a hard-coded
  fit or a claim of frame-perfect replay.

## Run

```sh
go run ./cmd/gcsim -c configs/stellar_swirl/current_account.gcsim -out current.json -nb
go run ./cmd/gcsim -c configs/stellar_swirl/prydwen_c1_kqms.gcsim -out public-c1.json -nb
go run ./cmd/gcsim -c configs/stellar_swirl/prydwen_c2_kqms.gcsim -out public-c2.json -nb
go run ./cmd/gcsim -c configs/stellar_swirl/public_168k_showcase.gcsim -out showcase.json -sampleMaxDps showcase-max.json -nb
```

## Sources

- KQM Stellar Glimmer mechanics and formulas:
  <https://keqingmains.com/misc/stellar-reaction-guide/>
- KQM Mizuki TCL: <https://library.keqingmains.com/characters/anemo/yumemizuki-mizuki>
- KQM Odette TCL: <https://library.keqingmains.com/characters/cryo/odette>
- KQM Cryo Traveler TCL: <https://library.keqingmains.com/characters/cryo/traveler-cryo>
- Exaiphanes Blade refinement data: <https://paimon.moe/weapons/exaiphanes_blade>
- Prydwen Mizuki builds, public gear, rotations, and relative comparisons:
  <https://www.prydwen.gg/genshin-impact/characters/yumemizuki-mizuki>
- Key_Influence_8972's recorded 168,551-DPS dummy test:
  <https://www.reddit.com/r/MizukiMainsGI/comments/1w3759x/168k_dps_c0_mizuki_dummy_test/>
- gcsim upstream: <https://github.com/genshinsim/gcsim>
