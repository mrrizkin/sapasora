# 11 — Pricing, Billing, dan Unit Economics

## Pricing principles

- Jangan menyembunyikan biaya provider/Meta.
- User dapat memperkirakan biaya sebelum send/campaign.
- Free tier harus aman dan terbatas, bukan channel produksi tanpa guardrail.
- Harga mengikuti value: inbox, automation, team, analytics, support.
- Overage memiliki hard limit/approval, bukan surprise bill.

## Suggested plans

### Sandbox — gratis

- Test/sandbox channel.
- Synthetic/test recipients.
- Developer Center.
- Limit rendah.
- Tidak untuk production campaign.

### Starter

- 1 workspace/channel.
- 3 agents.
- Basic inbox.
- Transactional notifications.
- Usage-based messaging.

### Growth

- Multiple agents/channels sesuai policy.
- Campaign approval.
- Automation.
- Templates.
- Analytics.
- Higher support level.

### Scale

- Advanced RBAC.
- Multiple WABA/channel.
- SLA/SLO agreement.
- Audit/export/retention controls.
- Dedicated support.

## Billing components

- Platform subscription.
- Provider/Meta usage cost.
- Media/storage overage.
- AI usage jika diaktifkan.
- Optional implementation/support fee.

Pisahkan invoice line item agar user paham apa yang dibayar ke platform dan apa yang berasal dari channel/provider.

## Billing rules

- Currency dan tax configuration.
- Monthly/annual commitment.
- Usage meter real-time dengan timestamp.
- Invoice preview sebelum renewal.
- Overage hard cap.
- Pause send jika limit tercapai, kecuali owner approve.
- Proration dijelaskan sebelum upgrade/downgrade.
- Refund/credit policy tertulis.
- Failed message dan retry dihitung sesuai definisi yang transparan.

## Unit economics

Track per workspace/channel:

```text
Gross revenue
- Meta/provider message cost
- media/storage cost
- compute/queue cost
- support cost
= contribution margin
```

Metrik:

- ARPA/ARPU.
- Gross margin.
- Cost per delivered message.
- Support cost per workspace.
- CAC payback.
- LTV/CAC.
- Free-to-paid conversion.
- Net revenue retention.

## Pricing experiments

1. Platform fee rendah + usage pass-through.
2. Paket berdasarkan agent/inbox + usage.
3. Paket berdasarkan automation/conversion value.
4. Annual discount dengan usage forecast.

Jangan mengunci harga sebelum provider cost, support cost, dan behavior customer pilot diketahui.

## Pricing page requirements

- Calculator 1.000/10.000/100.000 pesan.
- Contoh scenario transactional vs marketing.
- Effective monthly annual price.
- Channel/agent/template limits.
- Provider fee disclaimer.
- Overage estimate.
- Free sandbox CTA.
- Comparison table yang tidak bergantung pada toggle tersembunyi.
