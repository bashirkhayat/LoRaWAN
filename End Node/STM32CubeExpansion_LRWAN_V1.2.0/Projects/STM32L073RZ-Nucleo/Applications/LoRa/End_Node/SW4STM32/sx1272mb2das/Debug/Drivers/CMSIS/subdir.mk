################################################################################
# Automatically-generated file. Do not edit!
################################################################################

# Add inputs and outputs from these tool invocations to the build variables 
C_SRCS += \
/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/CMSIS/Device/ST/STM32L0xx/Source/Templates/system_stm32l0xx.c 

OBJS += \
./Drivers/CMSIS/system_stm32l0xx.o 

C_DEPS += \
./Drivers/CMSIS/system_stm32l0xx.d 


# Each subdirectory must supply rules for building sources it contributes
Drivers/CMSIS/system_stm32l0xx.o: /home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/CMSIS/Device/ST/STM32L0xx/Source/Templates/system_stm32l0xx.c
	@echo 'Building file: $<'
	@echo 'Invoking: MCU GCC Compiler'
	@echo $(PWD)
	arm-none-eabi-gcc -mcpu=cortex-m0plus -mthumb -mfloat-abi=soft -DSTM32L073xx -DUSE_STM32L0XX_NUCLEO -DUSE_HAL_DRIVER -DREGION_EU868 -DGXX_EXPERIMENTAL_CXX0X -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Projects/STM32L073RZ-Nucleo/Applications/LoRa/End_Node/Core/inc" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Projects/STM32L073RZ-Nucleo/Applications/LoRa/End_Node/LoRaWAN/App/inc" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/BSP/STM32L0xx_Nucleo" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/STM32L0xx_HAL_Driver/Inc" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/CMSIS/Device/ST/STM32L0xx/Include" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/CMSIS/Include" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Middlewares/Third_Party/LoRaWAN/Crypto" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Middlewares/Third_Party/LoRaWAN/Mac" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Middlewares/Third_Party/LoRaWAN/Phy" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Middlewares/Third_Party/LoRaWAN/Utilities" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Middlewares/Third_Party/LoRaWAN/Core" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/BSP/Components/Common" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/BSP/Components/hts221" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/BSP/Components/lps22hb" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/BSP/Components/lps25hb" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/BSP/Components/sx1272" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/BSP/X_NUCLEO_IKS01A1" -I"/home/ibraheem/Downloads/jeriess/STM32CubeExpansion_LRWAN_V1.2.0/Drivers/BSP/X_NUCLEO_IKS01A2" -I../../../../../../../../Drivers/BSP/sx1272mb2das  -Os -g3 -Wall -fmessage-length=0 -std=c++11 -ffunction-sections -c -fmessage-length=0 -MMD -MP -MF"$(@:%.o=%.d)" -MT"$@" -o "$@" "$<"
	@echo 'Finished building: $<'
	@echo ' '


