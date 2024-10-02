---
page_title: "Azure Marketplace"
subcategory: "Licensing"
---

# Azure Marketplace

## Pricing Model

* TerraProvider is offered as an **annual subscription plan** with different [User Segments](#user-segments). The correct **user segment** is automatically selected by our platform based on the amount of desired users.
* All subscription plans consist of a **base fee** which includes a certain amount of users per subscription cycle - depending on the **user segment**. For example, the **base fee** for the user segment _TerraProvider Enterprise (Y) 250_ includes 250 users per year.
* If more than the included amount of users is required, **additional users** can be added to the  plan. For each additional user, we charge an additional per-user fee.

## Invoicing

* During the first subscription interval, your subscription fees are not immediately due after completing the subscription enrolment. Instead we will start billing once your cancellation grace period has expired.&#x20;
* Upon every renewal date, you will be billed immediately.
* You will always be billed for the entire subscription cycle in advance.
* The related items should appear on your Microsoft Azure invoice (Pay-As-You-Go or Enterprise Agreement) the month after we have reported your fees to Microsoft.
*   In the PDF invoice you will receive from Microsoft, all TerraProvider fees are lumped into an item called "SaaS". The related Publisher is "glueckkanja".\


    <figure><img src="https://github.com/terraprovider/terraform-provider-microsoft365wp/raw/main/docs/.mdassets/image.png" alt="Screenshot" width="563"><figcaption></figcaption></figure>

-> For a more detailed cost breakdown of your base and additional user fees, please refer to the invoice in your Azure portal.

## Plan Overview

Subscriptions for TerraProvider are available based on an **annual** renewal interval.

| **Plan**                     | **Renewal Interval** |
| ---------------------------- | -------------------- |
| TerraProvider Enterprise (Y) | Annually             |

### User Segments

The following user segments are available:&#x20;

<table data-header-hidden><thead><tr><th width="240.02162801098973">Plan</th><th width="244.07580174927114">Included Users</th><th></th></tr></thead><tbody><tr><td><strong>User Segment</strong></td><td><strong>Included Users in Base Fee</strong></td><td><strong>Maximum Total Users</strong></td></tr><tr><td>TerraProvider Enterprise (Y) 100</td><td>50</td><td>249</td></tr><tr><td>TerraProvider Enterprise (Y) 250</td><td>250</td><td>999</td></tr><tr><td>TerraProvider Enterprise (Y) 1000</td><td>1,000</td><td>2,499</td></tr><tr><td>TerraProvider Enterprise (Y) 2500</td><td>2,500</td><td>9,999</td></tr><tr><td>TerraProvider Enterprise (Y) 10000</td><td>10,000</td><td>unlimited</td></tr></tbody></table>

For prices in Euro (EUR), please check out our [website](https://www.terraprovider.com/pricing/). For prices in _your_ currency, please directly refer to the **Marketplace** in the [Azure Portal](https://portal.azure.com/).

## Subscription Management

### User Upgrades

*   If you would like to upgrade your user count, you can do that any time during the current subscription cycle by navigating to your **TerraProvider subscription** in the [Azure SaaS portal](https://portal.azure.com/#blade/HubsExtension/BrowseResourceBlade/resourceType/Microsoft.SaaS%2Fresources) and by clicking "Open SaaS Account on publisher's site" (see screenshot below). This will re-direct you to our platform where the amount of users can be upgraded.\


    <figure><img src="https://github.com/terraprovider/terraform-provider-microsoft365wp/raw/main/docs/.mdassets/MASTER copy.jpg" alt="Screenshot"><figcaption></figcaption></figure>
* Our platform will inform you about the new fees you to expect for a **complete** subscription cycle.
* For the current cycle, we will bill the additional users for remaining days only.
* After confirming your choice and once we have updated the license in our backend, you will receive a confirmation email from us.

### User Downgrades

* Downgrading the amount of users is currently not possible without cancelling the subscription.
* If you want to perform a downgrade, please [cancel](#cancellation) your current subscription and re-subscribe with the desired amount of users once the cancellation becomes effective.

### **Recurring Billing**

If you decide to disable **Recurring billing**, your subscription will not renew automatically. Instead, Microsoft will (irreversibly) cancel the subscription towards the end of the current subscription cycle. This means, your license will be automatically revoked on that date as well. While the subscription has not expired yet, you can opt to enable **Recurring billing** at any time.

### Cancellation

*   If you would like to (irreversibly) cancel your subscription, navigate to your **RADIUSaaS subscription** in the [Azure SaaS portal](https://portal.azure.com/#blade/HubsExtension/BrowseResourceBlade/resourceType/Microsoft.SaaS%2Fresources) and click **Cancel subscription**.\


    <figure><img src="https://github.com/terraprovider/terraform-provider-microsoft365wp/raw/main/docs/.mdassets/MASTER copy (2).jpg" alt="Screenshot"><figcaption></figcaption></figure>
* If you cancel within the grace period, the license will be revoked immediately.
* If you cancel after the grace period, the license will remain active until the end of the current subscription cycle.

## **Trials**

In case you would like to test TerraProvider, please [get in contact with us](https://www.terraprovider.com/drop-a-question) or send us an email to [sales@terraprovider.com](mailto:sales@terraprovider.com).

## FAQs

### How to purchase TerraProvider?

To  get started with your TerraProvider subscription, follow below steps:

* Locate [TerraProvider](https://portal.azure.com/#view/Microsoft\_Azure\_Marketplace/GalleryItemDetailsBladeNopdl/id/glueckkanja-gabag.terraprovider-transactable-prod) on the **Marketplace** in your **Azure Portal**. In case we have extended a **Private Offer** to you or your MSP/distribution has extended a **Multiparty Offer (MPO)** to you, navigate to **Marketplace** in your **Azure Portal** and then to **Private Offer Management** to locate the Private Offer.
  * More details on Private Offers and MPOs can be found in Microsoft's documentation.
    * [Private Offer](https://learn.microsoft.com/en-us/marketplace/private-offers-purchase)
    * [Multiparty Offer](https://www.youtube.com/watch?v=TANUlgLuVqI)
*   Click **Subscribe**.\


    <figure><img src="https://github.com/terraprovider/terraform-provider-microsoft365wp/raw/main/docs/.mdassets/image (1).png" alt=""><figcaption></figcaption></figure>
* Create or select the **Resource group** you would like to deploy the subscription to.
* Assign a descriptive **Name** to later identify your subscription.
* We recommend to keep **Recurring billing** **On** so that you do not have to worry about an automatic termination of your subscription.
*   Click **Review + subscribe** and then **Subscribe** to deploy the **SaaS** resource to your **Resource group**.\


    <figure><img src="https://github.com/terraprovider/terraform-provider-microsoft365wp/raw/main/docs/.mdassets/image (2).png" alt=""><figcaption></figcaption></figure>

-> The random order of **Base Fees** und **Additional Users** under the **Price** information is attributed to limitations of the Azure Marketplace. Later during the the enrolment process, we will provide you with transparent information on the expected licensing fees.

*   Once the deployment is complete, please navigate to our platform to complete the checkout. Therefore click **Configure account now**.\


    <figure><img src="https://github.com/terraprovider/terraform-provider-microsoft365wp/raw/main/docs/.mdassets/image (3).png" alt=""><figcaption></figcaption></figure>
* After authenticating on our platform using your Microsoft credentials, you will be prompted for additional information, such as the desired total **User** amount and a **Technical contact**.

-> The **Technical contact** must have a mailbox connected to it, so we are able to notify you in case there are relevant issues with TerraProvider.

<figure><img src="https://github.com/terraprovider/terraform-provider-microsoft365wp/raw/main/docs/.mdassets/Screenshot_2024-09-13_at_15_25_38.jpg" alt=""><figcaption></figcaption></figure>

* Based on the amount of users provided, we will charge the relevant base fee for your user segment as well as additional users, in case you require more than the included amount in your base fee. **The platform automatically selects the best price / tier**.
* The platform will show you the licensing fees you have to expect under **Cost Projection**.
* If you are happy with it, please click **Review & Submit** for a final review and a fee summary.
* Complete the checkout by confirming your choice and clicking **Submit**.
* This triggers us to create a license for your organization. We will inform you via email with all relevant information on the next steps once this is completed. This won't take any longer than one business day.

-> You will only be charged by Microsoft, once you have completed the enrolment on our platform and received our welcome email.
