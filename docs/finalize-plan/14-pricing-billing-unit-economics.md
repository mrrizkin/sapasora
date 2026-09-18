# 14 — Pricing, Billing, dan Unit Economics

## Pricing principles

- Biaya platform dan provider/Meta dipisah.
- Estimasi cost sebelum send/campaign.
- Free tier sebagai sandbox aman, bukan unlimited production.
- Overage memiliki hard cap/approval.
- Value utama: inbox, automation, team, analytics, support.

## Suggested plans

### Sandbox — free

Test/sandbox channel, synthetic/test recipients, developer center, limit rendah, tanpa production campaign.

### Starter

1 workspace/channel, 3 agents, basic inbox, transactional notification, usage-based messaging.

### Growth

Multiple agents/channel sesuai policy, campaign approval, automation, templates, analytics, higher support.

### Scale

Advanced RBAC, multiple WABA/channel, audit/export/retention, SLA/SLO, dedicated support.

## Billing components

- Platform subscription.
- Meta/provider usage.
- Media/storage overage.
- AI usage.
- Implementation/support optional.

## Billing rules

- Currency/tax.
- Monthly/annual commitment.
- Real-time usage meter.
- Invoice preview.
- Overage hard cap.
- Pause send jika limit tercapai kecuali owner approve.
- Proration sebelum upgrade/downgrade.
- Refund/credit policy tertulis.
- Definisi biaya failed/retry message jelas.

## Unit economics

```text
Gross revenue
- Meta/provider message cost
- media/storage
- compute/queue
- support
= contribution margin
```

Track ARPA/ARPU, gross margin, cost per delivered message, support cost, CAC payback, LTV/CAC, free-to-paid, churn, NRR.

## Pricing experiments

- Platform fee + usage pass-through.
- Agent/inbox + usage.
- Automation/conversion value.
- Annual discount + usage forecast.

Jangan mengunci harga sebelum provider cost, support cost, dan pilot behavior diketahui.

## Pricing page

Wajib memiliki calculator 1K/10K/100K messages, transactional vs marketing scenario, annual effective price, channel/agent/template limits, provider fee disclaimer, overage estimate, sandbox CTA, dan comparison table tanpa hidden toggle.

## Discovery benchmark

Pricing Fonnte yang diamati menjadi benchmark, bukan target copy:

- Text Only monthly: Free Rp0; Lite Rp25K/1K; Regular Rp66K/10K; Regular Pro Rp110K/25K; Master Rp175K/unlimited.
- All Feature monthly: Super Rp165K/10K; Advanced Rp255K/25K; Ultra Rp355K/unlimited.
- Annual IDR: Lite Rp250K; Regular Rp660K; Regular Pro Rp1.1M; Master Rp1.75M; Super Rp1.65M; Advanced Rp2.55M; Ultra Rp3.55M.
- USD observed: monthly Lite `$2.5`, Regular `$6`, Regular Pro `$10`, Master `$14`, Super `$13`, Advanced `$22`, Ultra `$29`; annual Lite `$25`, Regular `$60`, Regular Pro `$100`, Master `$140`, Super `$130`, Advanced `$220`, Ultra `$290`.

Discovery juga menemukan paket per-device, 30 hari/bulan, 365 hari/tahun, subscribe ulang paket sama menambah quota/expiration, dan upgrade/downgrade dapat terminate paket lama. Produk kita harus menjelaskan aturan ini lebih transparan atau memilih model billing yang lebih sederhana.
