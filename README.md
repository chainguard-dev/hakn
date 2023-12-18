# `hakn`


`hakn` (pronounced like "kraken") stands for "High-Availability Knative", and is
an opinionated downstream distribution of Knative on which we (will) run
Chainguard's hosted platform.

### Who is this for?

We put this together to satisfy our own operational needs, and to give us a
place to codify some of the "opinionation" of our own internal platform.  Over
time, more of these opinions will surface.

### How is this different?

This is a "shrink wrapped" distribution of Knative, so where upstream offers a
wide selection of unopinionated "lego blocks" (e.g. networking layer, TLS
provider, eventing subsystem, ...), this is a mostly-built "lego set".  For the
most part, we will aim for upstream "conformance", but the underlying system may
itself behave fairly differently.

We operate our controlplane pods as HA-by-default stateful sets, which enables
us to bound the worst case downtime of a single replica being out of commission.
In our downstream CI, using upstream releases we would pretty frequently see
flakes due to one of the upstream webhooks getting restarted.  Each
`knative-serving` controlplane pod runs the following:
 - Upstream webhook logic (serving core, domain mappings)
 - Upstream controller logic (serving core, domain mappings, autoscaler)
 - Upstream `net-istio` integration (typically separate)
 - Custom GCLB TLS logic (see below, but based on Knative's webhook cert logic)

For our networking layer, we are bundling `net-istio` which is (currently) the
best-supported service mesh integration with Knative.  Currently we aren't
testing mesh-mode, but we plan to start soon; however, we also leverage Istio
(and a custom integration) to terminate TLS coming from Google's HTTPS load
balancer (GCLB).  Since we are terminating TLS at the GCLB and configuring a
self-signed cert for Istio to terminate at the `Gateway`, we do not utilize
Knative's auto-TLS integrations, since those fundamentally require SNI, which is
not supported by GCLB.

We do not (currently) support scaling to zero.  This is in part to simplify the
networking policies we have to author, since the way we have things configured
the upstream `activator` is never in the data path (we don't build or ship it).
To offset this, we have cranked up our default garbage collection settings to be
fairly aggressive about cleaning up old revisions.
