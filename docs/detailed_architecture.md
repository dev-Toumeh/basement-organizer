# Detailed Architecture


<!--toc:start-->
- [Detailed Architecture](#detailed-architecture)
  - [Architecture Overview](#architecture-overview)
  - [Functions](#functions)
  - [Place](#place)
  - [Item](#item)
<!--toc:end-->

## Architecture Overview
![](https://www.plantuml.com/plantuml/png/hLDVQzi-57tNfn2ce9SgEFwJsk77Zr9AHyksXz6oqZA6R2kE5baoaapJcdptMOeJEqa7PSmFZdLrphddt4YNsX1bcee3H3YNBnLRCr42JXS66WaLAvS65apHt30f03AojlLC6YEhqGOWkP4AUW-QAWyI3Hi0QfAlI4cXTttbIaK6Mm2xnvA8A7ctrNHuFVcbfN0GjxGhi1dnhxeFAWg7oBRV9qfSbamSrHGjcJPK7PLhelMBL8MsJHmFOZoqdlTQI4X3eGqi1UHCh210LgXjjxzBLb2uaD9aShsjr4GPbhEQ23EU13pTVI_I5hNP06JpsUO0x1B0iRyXrx53qFixiMSBw8_sFkJEePxn_xl2_72yiZPm-3DX8Qcs6I1TTnW0-wVdmf5ChYWnTCnrQ8lgD0-qdI0J8-m7XKyKd__3Gmwv0ltaXy7Rjqt_TGgqKKrk6aL1s-vxjllExWWBKb5ltdfDmmflxdKK9ahl77gxixOaUZbiRCvFwS5VyYUJISVSkCBs3GwIe7sK022w9bNDAK1z5EI6qPzhp_MU-JGCepp6Ypn8SEHV93XDiodsepXC4nhcSHfRPqDiZt2wtfdzsJiz2yZPD3gxZxCN_XXyVhr-MbrDQKFkN_sdvk4s_3QzkNdCIXxTVe_VxkvkllW-RTu-pHvp5jXywC1K58KBFmbZd0LaXgEeI74sIp8yJOCeoQAGfXV-iQcYOkBVkkelyW-0thkh-miA9nDxRj1X8G2N7RQf-6y0)

[Overview PlantUML](overview_plantuml.md)

## Functions
* Login
* Register
* View, create, modify item/box/shelf/area
* Delete box/shelf/area
    * What to do with items?
* View items in box/area/shelf
* Move items in/out of box/area/shelf
* Search for item/box/shelf/area
* Create QR code for item/box/(shelf?)
* Export data as JSON

![](https://www.plantuml.com/plantuml/png/TP11JiCm44NtESMNs41KAKTehm1RXD0G_pI3R4zbJvFszYXf4Ag0tRVythTd3jMbUD8SafoQ5iCvzwwMAvgsHRIZvZKL7hKwov8dvRn3NoZE7P857QvGP-e-xDBLaN5Y9uMom86HJlo-QTynZ-9mGp3CwkFQKVVVJzoEqBp8gxKrB42oYJz01pixRFYBh9JIZnYi_FV7ruShOxJvRPi9po_eBV3UURn9-wSxWvUJ5KSG5qZ5PxNS79Z3jj4l)

[Functions PlantUML](functions_plantuml.md)

## Place
* Room
    * Label
    * Description
    * Items
    * Boxes
    * shelves
    * shelf
    * Label
    * Description
    * Items
    * Rows/columns
    * Picture
    * QR code
    * Height/width/length?
    * Boxes
        * Label
        * Description
        * Items
        * Boxes
        * Picture
        * QR code

![](https://www.plantuml.com/plantuml/png/TP7FIW0X48VlFeNG8V2oZz2-GTGrzk3gi0hZ7tHiw-sRIe94FOtVuCptmxsIBjG22l5R1Xzjq5a8zVQ0_d9q5O1uo8pQW503BIa5XaoVYuuN71Br2IZ_CmlL59_9foZKtuWJ6_SyGQXCXy5d-e09hWxmdHkFYa-_yRAfmLfXwQ0uq5qtax25s4dsXbg1_keJ4DDzo9TNQP8TgOFleTljBRaD8LxadeHQXvjhhxtNulCDOu1PWXsYxP_-KxU0Nm00)

[Place PlantUML](place_plantuml.md)

## Item
* Label
* Description
* Picture
* Quantity
* Weight?
* QR code

![](https://www.plantuml.com/plantuml/png/LSqn3W8X58JXdbEGsEqbt1iOROt5MvYiBu6t24ECjnSI6xi_Np6pDafbLejPYblkfwJxarlgU7WbidVBm4iLsuSasP3-AQA5geLws71VDB1Np7nrCIe_iz-Oi-UPrviFHuHROV7t-mK0)

[Item PlantUML](item_plantuml.md)
